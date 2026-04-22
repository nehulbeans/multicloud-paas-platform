package models

import "time"

// DeploymentRequest is what the user sends in the POST /api/deploy body
type DeploymentRequest struct {
	AppName string            `json:"app_name"`
	Image   string            `json:"image"`
	Ports   []int32           `json:"ports"`
	EnvVars map[string]string `json:"env_vars"`

	// PaaS Specific Features
	Strategy     string   `json:"strategy"`                // "cost-optimized", "active-passive", "custom"
	CustomClouds []string `json:"custom_clouds,omitempty"` // Array of cloud IDs in priority order (used if strategy == "custom")
	Size         string   `json:"size"`                    // "micro" (512MB), "small" (1GB), "large" (4GB)
	Replicas     int32    `json:"replicas"`
}

// Deployment represents the logical application deployed by the user and its global routing state
type Deployment struct {
	ID             string    `json:"id" db:"id"`                             // UUID of the deployment
	AppName        string    `json:"app_name" db:"app_name"`                 // e.g., "test-nginx-app-20"
	Subdomain      string    `json:"subdomain" db:"subdomain"`               // e.g., "test-nginx-app-20.leancrust.dpdns.org"
	Strategy       string    `json:"strategy" db:"strategy"`                 // e.g., "active-passive"
	CurrentCloudID string    `json:"current_cloud_id" db:"current_cloud_id"` // The cloud currently receiving DNS traffic
	CreatedAt      time.Time `json:"created_at" db:"created_at"`

	// Targets holds the multi-cloud mappings (Populated dynamically via JOINs when needed)
	Targets []DeploymentCloud `json:"targets,omitempty" db:"-"`
}

// DeploymentCloud represents the presence and health of an app on a specific physical cloud
type DeploymentCloud struct {
	DeploymentID string `json:"deployment_id" db:"deployment_id"`
	CloudID      string `json:"cloud_id" db:"cloud_id"`
	Priority     int    `json:"priority" db:"priority"`     // 1 = Primary, 2 = Secondary, 3 = Tertiary
	AppStatus    string `json:"app_status" db:"app_status"` // "deploying", "running", "failed"
}
