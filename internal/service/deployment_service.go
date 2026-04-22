package service

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"

	"multicloud-paas-platform/internal/core/models"
	"multicloud-paas-platform/internal/repository"
	"multicloud-paas-platform/pkg/dnsclient"
	"multicloud-paas-platform/pkg/k8sclient"
)

type DeploymentService interface {
	DeployApp(ctx context.Context, req models.DeploymentRequest) (map[string]interface{}, error)
}

type deploymentService struct {
	clusterRepo    repository.ClusterRepository
	deploymentRepo repository.DeploymentRepository
	k8sDeployer    k8sclient.K8sDeployer
	dnsClient      *dnsclient.CloudflareClient
}

func NewDeploymentService(
	cluster_repo repository.ClusterRepository,
	deployment_repo repository.DeploymentRepository,
	deployer k8sclient.K8sDeployer,
	dns *dnsclient.CloudflareClient,
) DeploymentService {
	return &deploymentService{
		clusterRepo:    cluster_repo,
		deploymentRepo: deployment_repo,
		k8sDeployer:    deployer,
		dnsClient:      dns,
	}
}

func (s *deploymentService) DeployApp(ctx context.Context, req models.DeploymentRequest) (map[string]interface{}, error) {
	if req.AppName == "" || req.Image == "" {
		return nil, fmt.Errorf("app_name and image are required fields")
	}

	if req.Replicas == 0 {
		req.Replicas = 1
	}

	// 1. Generate the custom URL
	rootDomain := os.Getenv("ROOT_DOMAIN")
	if rootDomain == "" {
		rootDomain = "local.dev"
	}
	customURL := fmt.Sprintf("%s.%s", req.AppName, rootDomain)

	var targetClouds []models.Cloud

	// 2. Determine target clouds based on strategy
	if req.Strategy == "custom" {
		if len(req.CustomClouds) == 0 {
			return nil, fmt.Errorf("the custom_clouds array must contain at least one cloud ID")
		}

		// Double-check the requested clouds exist in our DB and fetch their Tunnel UUIDs
		for _, cloudID := range req.CustomClouds {
			cloud, err := s.clusterRepo.GetCloudByID(ctx, cloudID)
			if err != nil {
				fmt.Printf("Warning: Cloud %s not found or offline, skipping. %v\n", cloudID, err)
				continue
			}
			targetClouds = append(targetClouds, *cloud)
		}
	}

	if len(targetClouds) == 0 {
		return nil, fmt.Errorf("none of the requested clouds are currently available")
	}

	primaryCloud := targetClouds[0] // Priority 1

	// 3. Create the Deployment Record in the DB
	deploymentID := uuid.New().String()
	deploymentRecord := models.Deployment{
		ID:             deploymentID,
		AppName:        req.AppName,
		Subdomain:      customURL,
		Strategy:       req.Strategy,
		CurrentCloudID: primaryCloud.ID,
	}

	if err := s.deploymentRepo.CreateDeployment(ctx, deploymentRecord); err != nil {
		return nil, fmt.Errorf("failed to save deployment record: %v", err)
	}

	// 4. Deploy to ALL targeted clouds in the background
	for i, cloud := range targetClouds {
		priority := i + 1 // 1 = Primary, 2 = Secondary...

		// Fetch the decrypted Kubeconfig from the clusters table based on the cloud's Name
		kubeconfig, err := s.clusterRepo.GetKubeconfigByName(cloud.Name)
		if err != nil {
			fmt.Printf("Warning: Failed to retrieve kubeconfig for %s: %v\n", cloud.Name, err)
			s.deploymentRepo.AddDeploymentCloudMapping(ctx, models.DeploymentCloud{
				DeploymentID: deploymentID,
				CloudID:      cloud.ID,
				Priority:     priority,
				AppStatus:    "failed",
			})
			continue
		}

		// Push K8s manifests using the dynamically fetched Kubeconfig
		_, err = s.k8sDeployer.PushDeployment(ctx, kubeconfig, req, customURL)

		status := "running"
		if err != nil {
			fmt.Printf("Warning: Failed to deploy to %s: %v\n", cloud.Name, err)
			status = "failed"
		}

		// Map this specific cloud to the deployment in the DB
		s.deploymentRepo.AddDeploymentCloudMapping(ctx, models.DeploymentCloud{
			DeploymentID: deploymentID,
			CloudID:      cloud.ID,
			Priority:     priority,
			AppStatus:    status,
		})
	}

	// 5. DNS ROUTING - Point Cloudflare exactly to the Primary Cloud's Tunnel UUID
	if !strings.HasSuffix(customURL, "local.dev") {
		// MapToTunnel will use the CNAME record approach bypassing private IP limitations
		err := s.dnsClient.MapToTunnel(customURL, primaryCloud.TunnelUUID)
		if err != nil {
			fmt.Printf("Warning: DNS CNAME update failed: %v\n", err)
		} else {
			fmt.Printf("Successfully routed %s -> %s.cfargotunnel.com\n", customURL, primaryCloud.TunnelUUID)
		}
	}

	// 6. Return a clean, single response to the user
	return map[string]interface{}{
		"app_name":      req.AppName,
		"custom_url":    fmt.Sprintf("https://%s", customURL),
		"primary_cloud": primaryCloud.Name,
		"strategy":      req.Strategy,
		"status":        "deployed",
	}, nil
}
