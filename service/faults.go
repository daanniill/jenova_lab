package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

var injectedLatency atomic.Int64

func applyInjectedLatency() {
	delay := time.Duration(injectedLatency.Load())

	if delay > 0 {
		time.Sleep(delay)
	}
}

func latencyFaultHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	raw := r.URL.Query().Get("duration")

	delay, err := time.ParseDuration(raw)
	if err != nil {
		http.Error(w, "invalid duration", http.StatusBadRequest)
		return
	}

	if delay < 0 {
		http.Error(w, "duration cannot be negative", http.StatusBadRequest)
	}

	injectedLatency.Store(int64(delay))

	fmt.Fprintf(w, "latency set to %s\n", delay)
}
