package models

// DeploymentRequest is what the user sends in the POST /api/deploy body
type DeploymentRequest struct {
	AppName string            `json:"app_name"`
	Image   string            `json:"image"`
	Ports   []int32           `json:"ports"`
	EnvVars map[string]string `json:"env_vars"`

	// PaaS Specific Features
	Strategy    string `json:"strategy"`               // "cost-optimized", "high-availability", "custom"
	CustomCloud string `json:"custom_cloud,omitempty"` // Only used if strategy == "custom"
	Size        string `json:"size"`                   // "micro" (512MB), "small" (1GB), "large" (4GB)
	Replicas    int32  `json:"replicas"`
}

type Cluster struct {
	ID         string
	Name       string
	Kubeconfig string // encrypted
}
