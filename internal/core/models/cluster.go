package models

type Cluster struct {
	ID         string
	Name       string
	Kubeconfig string // encrypted
}

type AddClusterRequest struct {
	Name       string `json:"name"`       // e.g., "aws", "hetzner"
	Kubeconfig string `json:"kubeconfig"` // The raw, unencrypted kubeconfig string
}
