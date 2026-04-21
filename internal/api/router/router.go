package router

import (
	"net/http"

	"multicloud-paas-platform/internal/api/controllers"
)

// SetupRouter initializes the HTTP multiplexer and registers all API endpoints
func SetupRouter(
	deployController *controllers.DeploymentController,
	clusterController *controllers.ClusterController,
) *http.ServeMux {
	mux := http.NewServeMux()

	// for PaaS load balancers
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy"}`))
	})

	mux.HandleFunc("/api/deploy", func(w http.ResponseWriter, r *http.Request) {
		// cors headers.. setup correctly for production system
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		deployController.HandleDeploy(w, r)
	})

	mux.HandleFunc("/internal/api/clusters", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		clusterController.HandleAddCluster(w, r)
	})

	return mux
}
