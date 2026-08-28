package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
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
	
	http.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] received /work request", serviceName)

		result := WorkResponse{
			Service: serviceName,
			Status: "ok",
		}

		// Leaf service: nothing more to call
		if downstream == "" {
			writeJSON(w, http.StatusOK, result)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		// get the work endpoint of downstream url
		url := strings.TrimRight(downstream, "/") + "/work"

		// create new request to send to downstream service
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			result.Status = "error"
			result.Error = err.Error()
			writeJSON(w, http.StatusInternalServerError, result)
			return
		}

		// attempt to send request
		resp, err := client.Do(req)
		if err != nil {
			result.Status = "error"
			result.Error = err.Error()

			log.Printf("[%s] downstream error: %v", serviceName, err)

			writeJSON(w, http.StatusBadGateway, result)
			return
		}
		defer resp.Body.Close()

		// attempt to decode response 
		var child WorkResponse
		if err := json.NewDecoder(resp.Body).Decode(&child); err != nil {
			result.Status = "error"
			result.Error = fmt.Sprintf("invalid downstream response: %v", err)

			writeJSON(w, http.StatusBadGateway, result)
			return
		}

		result.Downstream = &child

		// checks for bad response
		if resp.StatusCode >= 400 {
			result.Status = "error"
			writeJSON(w, http.StatusBadGateway, result)
			return
		}

		writeJSON(w, http.StatusOK, result)
	})
	
	log.Printf("starting %s on : %s downstream=%q", serviceName, port, downstream,)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
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