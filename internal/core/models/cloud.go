package models

import "time"

type Cluster struct {
	ID         string
	Name       string
	Kubeconfig string // encrypted
}

type AddClusterRequest struct {
	Name       string `json:"name"`       // e.g., "aws", "hetzner"
	Kubeconfig string `json:"kubeconfig"` // The raw, unencrypted kubeconfig string
}

// CloudProfile is used by the Strategy Engine to weigh decisions
type CloudProfile struct {
	Name     string `json:"name"`
	CostTier int    `json:"cost_tier"`
}

// Cloud represents a physical K3s cluster and its associated Cloudflare Tunnel
type Cloud struct {
	ID         string    `json:"id" db:"id"`                   // e.g., "oracle-main", "aws-east", "laptop-local"
	Name       string    `json:"name" db:"name"`               // e.g., "Oracle Cloud (Frankfurt)"
	TunnelUUID string    `json:"tunnel_uuid" db:"tunnel_uuid"` // The Cloudflare Tunnel UUID for this cluster
	Status     string    `json:"status" db:"status"`           // "healthy", "down", "maintenance"
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
