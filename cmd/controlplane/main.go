package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"multicloud-paas-platform/internal/api/controllers"
	"multicloud-paas-platform/internal/api/router"
	"multicloud-paas-platform/internal/repository"
	"multicloud-paas-platform/internal/service"
	"multicloud-paas-platform/pkg/dnsclient"
	"multicloud-paas-platform/pkg/k8sclient"
)

func getDBConnectionString() string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// defaults if env vars are not set
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	if user == "" {
		user = "postgres"
	}
	if password == "" {
		password = "password"
	}
	if dbname == "" {
		dbname = "paas_db"
	}

	// postgres://user:password@host:port/dbname?sslmode=disable
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbname)
}

func main() {
	log.Println("Starting Multicloud PaaS Control Plane...")

	godotenv.Load()
	connStr := getDBConnectionString()

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Database is unreachable: %v", err)
	}
	log.Println("Connected to PostgreSQL database.")

	cfToken := os.Getenv("CLOUDFLARE_API_TOKEN")
	cfZoneID := os.Getenv("CLOUDFLARE_ZONE_ID")
	dnsClient := dnsclient.NewCloudflareClient(cfToken, cfZoneID)

	// repos
	clusterRepo := repository.NewClusterRepository(db)
	deploymentRepo := repository.NewDeploymentRepository(db)
	k8sDeployer := k8sclient.NewK8sDeployer()

	// services
	clusterService := service.NewClusterService(clusterRepo)
	deployService := service.NewDeploymentService(clusterRepo, deploymentRepo, k8sDeployer, dnsClient)

	// controllers
	clusterController := controllers.NewClusterController(clusterService)
	deployController := controllers.NewDeploymentController(deployService)

	mux := router.SetupRouter(deployController, clusterController)

	log.Println("Server listening on port 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
