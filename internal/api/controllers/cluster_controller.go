package controllers

import (
	"encoding/json"
	"net/http"

	"multicloud-paas-platform/internal/core/models"
	"multicloud-paas-platform/internal/service"
)

type ClusterController struct {
	clusterService service.ClusterService
}

func NewClusterController(s service.ClusterService) *ClusterController {
	return &ClusterController{clusterService: s}
}

func (c *ClusterController) HandleAddCluster(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.AddClusterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	err := c.clusterService.AddCluster(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Cluster registered and credentials encrypted securely",
	})
}
