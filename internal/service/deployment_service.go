package service

import (
	"context"
	"fmt"

	"multicloud-paas-platform/internal/core/models"
	"multicloud-paas-platform/internal/repository"
	"multicloud-paas-platform/pkg/k8sclient"
)

type DeploymentService interface {
	DeployApp(ctx context.Context, req models.DeploymentRequest) (map[string]string, error)
}

type deploymentService struct {
	clusterRepo repository.ClusterRepository
	k8sDeployer k8sclient.K8sDeployer
}

func NewDeploymentService(repo repository.ClusterRepository, deployer k8sclient.K8sDeployer) DeploymentService {
	return &deploymentService{
		clusterRepo: repo,
		k8sDeployer: deployer,
	}
}

func (s *deploymentService) DeployApp(ctx context.Context, req models.DeploymentRequest) (map[string]string, error) {
	if req.AppName == "" || req.Image == "" {
		return nil, fmt.Errorf("app_name and image are required fields")
	}

	// default at least 1 replica
	if req.Replicas == 0 {
		req.Replicas = 1
	}

	// available clouds and their cost tiers
	// TODO: store and fetch this dynamically from db
	availableClouds := []models.CloudProfile{
		{Name: "hetzner", CostTier: 1}, // Cheap
		{Name: "aws", CostTier: 3},     // Expensive
	}

	strategy := GetStrategy(req.Strategy)
	deploymentPlan := strategy.GeneratePlan(req, availableClouds)

	// execute deployments and collect endpoints
	deploymentResults := make(map[string]string)

	for cloudName, tailoredReq := range deploymentPlan {
		kubeconfig, err := s.clusterRepo.GetKubeconfigByName(cloudName)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve credentials for %s: %v", cloudName, err)
		}

		publicEndpoint, err := s.k8sDeployer.PushDeployment(ctx, kubeconfig, tailoredReq)
		if err != nil {
			return nil, fmt.Errorf("deployment to %s failed: %v", cloudName, err)
		}

		deploymentResults[cloudName] = publicEndpoint
	}

	return deploymentResults, nil
}
