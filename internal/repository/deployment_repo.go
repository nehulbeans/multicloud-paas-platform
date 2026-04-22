package repository

import (
	"context"
	"database/sql"
	"fmt"

	"multicloud-paas-platform/internal/core/models"
)

// DeploymentRepository handles all database operations for apps and failover states
type DeploymentRepository interface {
	// Deployment Creation
	CreateDeployment(ctx context.Context, d models.Deployment) error
	AddDeploymentCloudMapping(ctx context.Context, dc models.DeploymentCloud) error
	// Failover & Health Worker Functions
	UpdateCurrentCloud(ctx context.Context, deploymentID, newCloudID string) error
	UpdateAppStatus(ctx context.Context, deploymentID, cloudID, status string) error
	GetNextHealthyCloud(ctx context.Context, deploymentID string) (*models.Cloud, error)
}

type deploymentRepository struct {
	db *sql.DB
}

func NewDeploymentRepository(db *sql.DB) DeploymentRepository {
	return &deploymentRepository{db: db}
}

// CreateDeployment inserts the main application record
func (r *deploymentRepository) CreateDeployment(ctx context.Context, d models.Deployment) error {
	query := `
		INSERT INTO deployments (id, app_name, subdomain, strategy, current_cloud_id)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query, d.ID, d.AppName, d.Subdomain, d.Strategy, d.CurrentCloudID)
	if err != nil {
		return fmt.Errorf("failed to insert deployment: %v", err)
	}
	return nil
}

// AddDeploymentCloudMapping links the deployment to a specific cloud with a priority rank
func (r *deploymentRepository) AddDeploymentCloudMapping(ctx context.Context, dc models.DeploymentCloud) error {
	query := `
		INSERT INTO deployment_clouds (deployment_id, cloud_id, priority, app_status)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(ctx, query, dc.DeploymentID, dc.CloudID, dc.Priority, dc.AppStatus)
	if err != nil {
		return fmt.Errorf("failed to insert deployment_cloud mapping: %v", err)
	}
	return nil
}

// =====================================================================
// HEALTH WORKER & FAILOVER FUNCTIONS
// =====================================================================

// UpdateCurrentCloud is called when a failover happens to record where DNS is currently pointing
func (r *deploymentRepository) UpdateCurrentCloud(ctx context.Context, deploymentID, newCloudID string) error {
	query := `UPDATE deployments SET current_cloud_id = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, newCloudID, deploymentID)
	if err != nil {
		return fmt.Errorf("failed to update current_cloud_id: %v", err)
	}
	return nil
}

// UpdateAppStatus is called by the health checker when a specific K3s cluster drops offline
func (r *deploymentRepository) UpdateAppStatus(ctx context.Context, deploymentID, cloudID, status string) error {
	query := `UPDATE deployment_clouds SET app_status = $1 WHERE deployment_id = $2 AND cloud_id = $3`
	_, err := r.db.ExecContext(ctx, query, status, deploymentID, cloudID)
	if err != nil {
		return fmt.Errorf("failed to update app status: %v", err)
	}
	return nil
}

// GetNextHealthyCloud finds the lowest-priority-number cloud that is still marked as 'running'
func (r *deploymentRepository) GetNextHealthyCloud(ctx context.Context, deploymentID string) (*models.Cloud, error) {
	query := `
		SELECT c.id, c.name, c.tunnel_uuid, c.status
		FROM clouds c
		JOIN deployment_clouds dc ON c.id = dc.cloud_id
		WHERE dc.deployment_id = $1
		  AND dc.app_status = 'running'
		  AND c.status = 'healthy'
		ORDER BY dc.priority ASC
		LIMIT 1
	`

	row := r.db.QueryRowContext(ctx, query, deploymentID)

	var cloud models.Cloud
	err := row.Scan(&cloud.ID, &cloud.Name, &cloud.TunnelUUID, &cloud.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no healthy backup clouds available for deployment %s", deploymentID)
		}
		return nil, fmt.Errorf("failed to query next healthy cloud: %v", err)
	}

	return &cloud, nil
}
