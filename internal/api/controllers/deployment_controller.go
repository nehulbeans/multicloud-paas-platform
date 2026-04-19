package controllers

import (
	"encoding/json"
	"net/http"

	"multicloud-paas-platform/internal/core/models"
	"multicloud-paas-platform/internal/service"
)

type DeploymentController struct {
	deployService service.DeploymentService
}

func NewDeploymentController(s service.DeploymentService) *DeploymentController {
	return &DeploymentController{deployService: s}
}

func (c *DeploymentController) HandleDeploy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.DeploymentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// this will block until IPs are provisioned
	results, err := c.deployService.DeployApp(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := map[string]interface{}{
		"status":    "success",
		"message":   "Deployment completed successfully",
		"endpoints": results,
	}

	json.NewEncoder(w).Encode(response)
}
