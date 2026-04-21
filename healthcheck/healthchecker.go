package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)


var (
	primary       = getEnv("PRIMARY_URL", "http://localhost:5001/healthz") //primary cluster url
	secondary     = getEnv("SECONDARY_URL", "http://localhost:5002/healthz") //secondary cluster url
	checkInterval = getDuration("CHECK_INTERVAL", 5*time.Second)
	httpTimeout   = getDuration("HTTP_TIMEOUT", 2*time.Second)
)

var currentActive = "PRIMARY"

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}

func isAlive(url string) bool {
	client := http.Client{Timeout: httpTimeout}
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

func main() {
	fmt.Printf("Starting healthchecker\n  primary:   %s\n  secondary: %s\n  interval:  %s\n",
		primary, secondary, checkInterval)

	for {
		primaryAlive := isAlive(primary)

		if !primaryAlive && currentActive == "PRIMARY" {
			fmt.Printf("[%s] ⚠️  Primary DOWN → Switching to SECONDARY\n", time.Now().Format(time.RFC3339))
			currentActive = "SECONDARY"
		} else if primaryAlive && currentActive == "SECONDARY" {
			fmt.Printf("[%s] ✅ Primary BACK → Switching to PRIMARY\n", time.Now().Format(time.RFC3339))
			currentActive = "PRIMARY"
		} else {
			fmt.Printf("[%s] Active: %s | primary alive: %v\n",
				time.Now().Format(time.RFC3339), currentActive, primaryAlive)
		}

		time.Sleep(checkInterval)
	}
}