package service

import (
	"fmt"

	"multicloud-paas-platform/internal/core/models"
	"multicloud-paas-platform/internal/repository"
	"multicloud-paas-platform/pkg/crypto"
)

type ClusterService interface {
	AddCluster(req models.AddClusterRequest) error
}

type clusterService struct {
	clusterRepo repository.ClusterRepository
}

func NewClusterService(repo repository.ClusterRepository) ClusterService {
	return &clusterService{
		clusterRepo: repo,
	}
}

func (s *clusterService) AddCluster(req models.AddClusterRequest) error {
	if req.Name == "" || req.Kubeconfig == "" {
		return fmt.Errorf("cluster name and kubeconfig are required")
	}

	// encrypt the raw kubeconfig using your AES-256-GCM function
	encryptedKubeconfig, err := crypto.Encrypt(req.Kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to encrypt kubeconfig: %v", err)
	}

	err = s.clusterRepo.AddCluster(req.Name, encryptedKubeconfig)
	if err != nil {
		return err
	}

	return nil
}
