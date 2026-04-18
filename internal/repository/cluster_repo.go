package repository

import (
	"database/sql"
	"fmt"
)

// ClusterRepository defines data operations for clusters
type ClusterRepository interface {
	GetKubeconfigByName(cloudName string) (string, error)
}

type clusterRepo struct {
	db *sql.DB
}

// NewClusterRepository creates a new instance of ClusterRepository
func NewClusterRepository(db *sql.DB) ClusterRepository {
	return &clusterRepo{db: db}
}

// GetKubeconfigByName fetches the raw kubeconfig string from the database
// will have to encrypt/decrypt the kubeconfig before returning: TODO
func (r *clusterRepo) GetKubeconfigByName(cloudName string) (string, error) {
	var kubeconfig string
	query := "SELECT kubeconfig FROM clusters WHERE name = $1"

	err := r.db.QueryRow(query, cloudName).Scan(&kubeconfig)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("cluster '%s' not found in database", cloudName)
		}
		return "", fmt.Errorf("database error: %v", err)
	}
	return kubeconfig, nil
}
