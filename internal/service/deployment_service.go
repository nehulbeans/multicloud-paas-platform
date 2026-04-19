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

	// deployment algorithm: Todo/ partially done
	var targetClouds []string
	switch req.Strategy {
	case "cost-optimized":
		targetClouds = []string{"ocs"}
	case "high-availability":
		targetClouds = []string{"aws", "ocs"}
	case "custom":
		if req.CustomCloud != "" {
			targetClouds = []string{req.CustomCloud}
		} else {
			return nil, fmt.Errorf("custom_cloud must be specified when strategy is custom")
		}
	default:
		targetClouds = []string{"ocs"}
	}

	// execute deployments and collect endpoints
	deploymentResults := make(map[string]string)

	for _, cloudName := range targetClouds {
		kubeconfig, err := s.clusterRepo.GetKubeconfigByName(cloudName)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve credentials for %s: %v", cloudName, err)
		}

		// push to k8s cluster and wait for the loadBalancer IP
		publicEndpoint, err := s.k8sDeployer.PushDeployment(ctx, kubeconfig, req)
		if err != nil {
			return nil, fmt.Errorf("deployment to %s failed: %v", cloudName, err)
		}

		deploymentResults[cloudName] = publicEndpoint
	}

	return deploymentResults, nil
}
