package main

import (
	"net/http"
	"os"
	"time"
)

type WorkResponse struct {
	Service    string        `json:"service"`
	Status     string        `json:"status"`
	Downstream *WorkResponse `json:"downstream,omitempty"`
	Error      string        `json:"error,omitempty"`
}

func main() {
	serviceName := getEnv("SERVICE_NAME", "service")
	port := getEnv("PORT", "8001")
	downstream := os.Getenv("DOWNSTREAM_URL")

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}