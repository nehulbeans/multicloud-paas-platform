package service

import (
	"multicloud-paas-platform/internal/core/models"
	"sort"
	"strings"
)

// Strategy Interface
type DeploymentStrategy interface {
	// GeneratePlan takes the user request and a list of available clouds,
	// and returns a deployment plan mapped by cloud name.
	GeneratePlan(req models.DeploymentRequest, availableClouds []models.CloudProfile) map[string]models.DeploymentRequest
}

func GetStrategy(strategyName string) DeploymentStrategy {
	switch strings.ToLower(strategyName) {
	case "scale-to-zero":
		return &ColdStandbyStrategy{}
	case "pilot-light":
		return &WarmStandbyStrategy{}
	case "active-active":
		return &ActiveActiveStrategy{}
	default:
		return &ColdStandbyStrategy{} // Default fallback
	}
}

// ---------------------------------------------------------
// Helper Function: Sort clouds by Cost (Cheapest first)
// ---------------------------------------------------------
func sortCloudsByCost(clouds []models.CloudProfile) {
	sort.Slice(clouds, func(i, j int) bool {
		return clouds[i].CostTier < clouds[j].CostTier
	})
}

// ---------------------------------------------------------
// ALGORITHM 1: Cold Standby (Scale-to-Zero)
// Logic: 100% replicas on the cheapest cloud. 0 on all others.
// ---------------------------------------------------------
type ColdStandbyStrategy struct{}

func (s *ColdStandbyStrategy) GeneratePlan(req models.DeploymentRequest, availableClouds []models.CloudProfile) map[string]models.DeploymentRequest {
	plan := make(map[string]models.DeploymentRequest)
	if len(availableClouds) == 0 {
		return plan
	}

	sortCloudsByCost(availableClouds)
	cheapestCloud := availableClouds[0].Name

	for _, cloud := range availableClouds {
		cloudReq := req
		if cloud.Name == cheapestCloud {
			cloudReq.Replicas = req.Replicas // Full compute
		} else {
			cloudReq.Replicas = 0 // Zero compute
		}
		plan[cloud.Name] = cloudReq
	}

	return plan
}

// ---------------------------------------------------------
// ALGORITHM 2: Warm Standby (Pilot Light)
// Logic: Full replicas on the cheapest. Exactly 1 replica on the rest.
// ---------------------------------------------------------
type WarmStandbyStrategy struct{}

func (s *WarmStandbyStrategy) GeneratePlan(req models.DeploymentRequest, availableClouds []models.CloudProfile) map[string]models.DeploymentRequest {
	plan := make(map[string]models.DeploymentRequest)
	if len(availableClouds) == 0 {
		return plan
	}

	sortCloudsByCost(availableClouds)
	cheapestCloud := availableClouds[0].Name

	for _, cloud := range availableClouds {
		cloudReq := req
		if cloud.Name == cheapestCloud {
			cloudReq.Replicas = req.Replicas // Full compute
		} else {
			cloudReq.Replicas = 1 // Pilot Light (1 pod)
		}
		plan[cloud.Name] = cloudReq
	}

	return plan
}

// ---------------------------------------------------------
// ALGORITHM 3: Active-Active (Max HA)
// Logic: User's requested replicas are deployed exactly as requested to ALL clouds.
// ---------------------------------------------------------
type ActiveActiveStrategy struct{}

func (s *ActiveActiveStrategy) GeneratePlan(req models.DeploymentRequest, availableClouds []models.CloudProfile) map[string]models.DeploymentRequest {
	plan := make(map[string]models.DeploymentRequest)

	for _, cloud := range availableClouds {
		cloudReq := req
		cloudReq.Replicas = req.Replicas // 100% redundancy everywhere
		plan[cloud.Name] = cloudReq
	}

	return plan
}
