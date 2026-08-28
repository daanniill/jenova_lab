package main

import (
	"encoding/json"
	"log"
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

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request){
		writeJSON(w, http.StatusOK, map[string]string {
			"service": serviceName,
			"status": "ok",
		})
	}) 
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("failed to encode respone: %v", err)
	}
}