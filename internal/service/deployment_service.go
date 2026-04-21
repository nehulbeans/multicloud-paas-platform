package service

import (
	"context"
	"fmt"
	"os"
	"strings"

	"multicloud-paas-platform/internal/core/models"
	"multicloud-paas-platform/internal/repository"
	"multicloud-paas-platform/pkg/dnsclient"
	"multicloud-paas-platform/pkg/k8sclient"
)

type DeploymentService interface {
	DeployApp(ctx context.Context, req models.DeploymentRequest) (map[string]string, error)
}

type deploymentService struct {
	clusterRepo repository.ClusterRepository
	k8sDeployer k8sclient.K8sDeployer
	dnsClient   *dnsclient.CloudflareClient
}

// NewDeploymentService injects the Database, Kubernetes, and DNS clients
func NewDeploymentService(
	repo repository.ClusterRepository,
	deployer k8sclient.K8sDeployer,
	dns *dnsclient.CloudflareClient,
) DeploymentService {
	return &deploymentService{
		clusterRepo: repo,
		k8sDeployer: deployer,
		dnsClient:   dns,
	}
}

// this orchestrates the multi-cloud deployment and configures global DNS
func (s *deploymentService) DeployApp(ctx context.Context, req models.DeploymentRequest) (map[string]string, error) {
	if req.AppName == "" || req.Image == "" {
		return nil, fmt.Errorf("app_name and image are required fields")
	}

	// default at least 1 replica
	if req.Replicas == 0 {
		req.Replicas = 1
	}

	// Step 1: Generate the custom URL FIRST
	// We must do this before deploying so we can tell the Kubernetes Ingress what domain to listen for
	rootDomain := os.Getenv("ROOT_DOMAIN") // e.g., "leancrust.dpdns.org"
	if rootDomain == "" {
		rootDomain = "local.dev" // Fallback for local testing
	}
	customURL := fmt.Sprintf("%s.%s", req.AppName, rootDomain)

	// available clouds and their cost tiers
	// TODO: store and fetch this dynamically from db
	availableClouds := []models.CloudProfile{
		{Name: "oracle", CostTier: 1}, // Cheap
	}

	// gets the algorithmic strategy and generate the plan
	strategy := GetStrategy(req.Strategy)
	deploymentPlan := strategy.GeneratePlan(req, availableClouds)

	// execute deployments and collect endpoints
	deploymentResults := make(map[string]string)

	for cloudName, tailoredReq := range deploymentPlan {
		kubeconfig, err := s.clusterRepo.GetKubeconfigByName(cloudName)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve credentials for %s: %v", cloudName, err)
		}

		// Step 2: Push the tailored k8s objects AND pass the customURL to create the Ingress
		publicEndpoint, err := s.k8sDeployer.PushDeployment(ctx, kubeconfig, tailoredReq, customURL)
		if err != nil {
			return nil, fmt.Errorf("deployment to %s failed: %v", cloudName, err)
		}

		deploymentResults[cloudName] = publicEndpoint
	}

	// Step 3: DNS ROUTING - map the custom domain to the primary active cloud IP
	_, err := s.configureGlobalDNS(customURL, req.Strategy, deploymentResults)
	if err != nil {
		// we log the error but don't fail the deployment since the K8s pods are running
		fmt.Printf("Warning: DNS update failed: %v\n", err)
	}

	// add the custom URL to the results so the frontend can display a clickable link to the user
	deploymentResults["custom_url"] = customURL

	return deploymentResults, nil
}

// configureGlobalDNS assigns the application a custom URL and maps it to the active cloud IP
func (s *deploymentService) configureGlobalDNS(customURL string, strategy string, deploymentResults map[string]string) (string, error) {
	// Skip Cloudflare API call if running locally
	if strings.HasSuffix(customURL, "local.dev") {
		return customURL, nil
	}

	// Determine the Primary IP to point the DNS to
	var primaryIP string
	primaryIP = deploymentResults["oracle"]

	if strings.HasPrefix(primaryIP, "192.168.") {
		fmt.Printf("Skipping DNS update: %s is behind a Cloudflare Tunnel\n", primaryIP)
		return customURL, nil
	}

	if primaryIP == "" {
		return customURL, fmt.Errorf("primary IP is empty, cannot update DNS")
	}

	// Strip the port to satisfy Cloudflare API rules
	cleanIP := strings.Split(primaryIP, ":")[0]

	// Send the Upsert request to Cloudflare
	err := s.dnsClient.UpsertRecord(customURL, cleanIP)
	if err != nil {
		return customURL, fmt.Errorf("failed to map %s to %s: %v", customURL, cleanIP, err)
	}

	fmt.Printf("Successfully mapped %s to %s\n", customURL, cleanIP)
	return customURL, nil
}
