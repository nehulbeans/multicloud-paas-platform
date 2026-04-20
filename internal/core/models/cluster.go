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

type CloudProfile struct {
	Name     string
	CostTier int // e.g., 1 = Cheap (Hetzner), 2 = Medium (DigitalOcean), 3 = Expensive (AWS/GCP)
}
