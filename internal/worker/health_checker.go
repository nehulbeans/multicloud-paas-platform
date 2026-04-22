package worker

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"multicloud-paas-platform/internal/core/models"
	"multicloud-paas-platform/internal/repository"
	"multicloud-paas-platform/pkg/dnsclient"
)

type HealthCheckWorker struct {
	deployRepo repository.DeploymentRepository
	dnsClient  *dnsclient.CloudflareClient
	interval   time.Duration
	client     *http.Client
}

// NewHealthCheckWorker initializes the background worker
func NewHealthCheckWorker(repo repository.DeploymentRepository, dns *dnsclient.CloudflareClient) *HealthCheckWorker {
	return &HealthCheckWorker{
		deployRepo: repo,
		dnsClient:  dns,
		interval:   30 * time.Second, // Check every 30 seconds
		client: &http.Client{
			Timeout: 5 * time.Second, // Don't wait forever for a dead app
			// Do not follow redirects for health checks (so we don't get tricked by auth walls)
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

// Start begins the infinite loop of health checking in a background goroutine
func (w *HealthCheckWorker) Start(ctx context.Context) {
	log.Println("Starting Multi-Cloud Health Check Worker...")
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Health Check Worker shutting down.")
			return
		case <-ticker.C:
			w.checkAllDeployments(ctx)
		}
	}
}

func (w *HealthCheckWorker) checkAllDeployments(ctx context.Context) {
	deployments, err := w.deployRepo.GetAllActiveDeployments(ctx)
	if err != nil {
		log.Printf("Worker Error: Failed to fetch deployments: %v", err)
		return
	}

	// 1. Tell us how many it found!
	log.Printf("🔍 Health Worker Woke Up: Found %d active deployments in DB.", len(deployments))

	for _, dep := range deployments {
		// 2. Fix Timezone mismatch (convert negative elapsed time to positive)
		elapsed := time.Since(dep.CreatedAt)
		if elapsed < 0 {
			elapsed = -elapsed
		}

		if elapsed < 45*time.Second {
			log.Printf("⏳ Skipping %s (Too new: deployed %v ago)", dep.AppName, elapsed)
			continue
		}

		// 3. Log that it is actively testing the URL
		log.Printf("🌐 Pinging: https://%s ...", dep.Subdomain)
		go w.checkSingleDeployment(ctx, dep)
	}
}

func (w *HealthCheckWorker) checkSingleDeployment(ctx context.Context, dep models.Deployment) {
	url := fmt.Sprintf("https://%s", dep.Subdomain)

	// Perform the health check with retries
	isHealthy := w.pingWithRetries(url, 3)

	if isHealthy {
		log.Printf("✅ %s is HEALTHY.", dep.AppName)
		return
	}

	log.Printf("🚨 FAILOVER ALERT: App %s (%s) failed 3 consecutive health checks!", dep.AppName, url)

	// ==========================================
	// EXECUTE FAILOVER
	// ==========================================

	// 1. Mark the current cloud as failed in the DB
	err := w.deployRepo.UpdateAppStatus(ctx, dep.ID, dep.CurrentCloudID, "failed")
	if err != nil {
		log.Printf("Worker Error: Failed to update app status for %s: %v", dep.AppName, err)
	}

	// 2. Ask the DB for the next healthy cloud based on priority
	backupCloud, err := w.deployRepo.GetNextHealthyCloud(ctx, dep.ID)
	if err != nil {
		log.Printf("💀 FATAL: All clouds are down for app %s! No failover possible.", dep.AppName)
		return
	}

	log.Printf("🔄 Rerouting %s to backup cloud: %s", dep.AppName, backupCloud.Name)

	// 3. Re-route internet traffic instantly to the backup cloud's Tunnel UUID
	// (We skip Cloudflare API if using local testing domains)
	if !strings.HasSuffix(dep.Subdomain, "local.dev") {
		err = w.dnsClient.MapToTunnel(dep.Subdomain, backupCloud.TunnelUUID)
		if err != nil {
			log.Printf("Worker Error: Failed to update Cloudflare DNS for %s: %v", dep.AppName, err)
			return
		}
	} else {
		log.Printf("⚠️ Skipped Cloudflare API because domain ends in local.dev")
	}

	// 4. Record the successful failover in the database
	err = w.deployRepo.UpdateCurrentCloud(ctx, dep.ID, backupCloud.ID)
	if err != nil {
		log.Printf("Worker Error: Failed to update current_cloud_id for %s: %v", dep.AppName, err)
	}

	log.Printf("✅ Failover complete for %s. Traffic is now flowing to %s.", dep.AppName, backupCloud.Name)
}

// pingWithRetries attempts to GET the URL. If it fails or returns a bad status code, it retries up to maxRetries.
func (w *HealthCheckWorker) pingWithRetries(url string, maxRetries int) bool {
	for attempt := 1; attempt <= maxRetries; attempt++ {
		req, _ := http.NewRequest("GET", url, nil)
		// Add a custom User-Agent so users see our platform pinging them, not a random bot
		req.Header.Set("User-Agent", "LeanCrust-PaaS-HealthCheck/1.0")

		resp, err := w.client.Do(req)

		// Check if it was a complete network failure or a bad status code
		if err != nil || resp.StatusCode >= 500 {
			if resp != nil {
				resp.Body.Close()
			}

			// If this isn't the last attempt, wait 3 seconds and try again
			if attempt < maxRetries {
				time.Sleep(3 * time.Second)
				continue
			}
			return false // All retries failed
		}

		resp.Body.Close()

		// If we get a 200-499, the cluster is alive and routing traffic.
		// (Even a 401 Unauthorized or 404 from the user's specific app means K3s and Traefik are healthy!)
		return true
	}
	return false
}
