package repository

import (
	"database/sql"
	"fmt"
	"multicloud-paas-platform/pkg/crypto"
)

// ClusterRepository defines data operations for clusters
type ClusterRepository interface {
	AddCluster(name string, encryptedKubeconfig string) error
	GetKubeconfigByName(cloudName string) (string, error)
}

type clusterRepo struct {
	db *sql.DB
}

// NewClusterRepository creates a new instance of ClusterRepository
func NewClusterRepository(db *sql.DB) ClusterRepository {
	return &clusterRepo{db: db}
}

func (r *clusterRepo) AddCluster(name string, encryptedKubeconfig string) error {
	query := "INSERT INTO clusters (name, kubeconfig) VALUES ($1, $2)"

	_, err := r.db.Exec(query, name, encryptedKubeconfig)
	if err != nil {
		return fmt.Errorf("failed to insert cluster into database: %v", err)
	}

	return nil
}

// GetKubeconfigByName fetches the raw kubeconfig string from the database
func (r *clusterRepo) GetKubeconfigByName(cloudName string) (string, error) {
	var encryptedKubeconfig string
	query := "SELECT kubeconfig FROM clusters WHERE name = $1"

	err := r.db.QueryRow(query, cloudName).Scan(&encryptedKubeconfig)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("cluster '%s' not found in database", cloudName)
		}
		return "", fmt.Errorf("database error: %v", err)
	}

	// Decrypt it before handing it to the Kubernetes Deployer
	decryptedKubeconfig, err := crypto.Decrypt(encryptedKubeconfig)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt kubeconfig: %v", err)
	}

	return decryptedKubeconfig, nil
}
