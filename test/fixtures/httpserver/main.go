package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	port := getEnv("PORT", "43102")
	exitImmediate := getEnv("EXIT_IMMEDIATELY", "") != ""

	if exitImmediate {
		os.Exit(1)
	}

	startupDelay := getEnvInt("STARTUP_DELAY_MS", 0)
	if startupDelay > 0 {
		time.Sleep(time.Duration(startupDelay) * time.Millisecond)
	}

	healthStatus := getEnvInt("HEALTH_STATUS", 200)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(healthStatus)
		fmt.Fprintf(w, `{"status": "%s"}`, statusText(healthStatus))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	addr := "127.0.0.1:" + port
	fmt.Fprintf(os.Stderr, "fixture server listening on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func statusText(code int) string {
	switch code {
	case 200:
		return "healthy"
	case 503:
		return "unavailable"
	default:
		return "unknown"
	}
}
