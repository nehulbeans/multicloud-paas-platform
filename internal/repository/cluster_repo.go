package repository

import (
	"context"
	"database/sql"
	"fmt"
	"multicloud-paas-platform/internal/core/models"
	"multicloud-paas-platform/pkg/crypto"
)

// ClusterRepository defines data operations for clusters
type ClusterRepository interface {
	AddCluster(name string, encryptedKubeconfig string) error
	GetKubeconfigByName(cloudName string) (string, error)
	GetCloudByID(ctx context.Context, id string) (*models.Cloud, error) // NEW
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

func (r *clusterRepo) GetCloudByID(ctx context.Context, id string) (*models.Cloud, error) {
	var cloud models.Cloud
	query := `SELECT id, name, tunnel_uuid, status, created_at FROM clouds WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&cloud.ID,
		&cloud.Name,
		&cloud.TunnelUUID,
		&cloud.Status,
		&cloud.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("cloud with ID '%s' not found", id)
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	return &cloud, nil
}
