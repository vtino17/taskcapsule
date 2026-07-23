package health

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"syscall"
	"time"
)

type Config struct {
	Type           string `json:"type"`
	URL            string `json:"url"`
	Host           string `json:"host"`
	Port           string `json:"port"`
	ExpectedStatus int    `json:"expectedStatus"`
	TimeoutSeconds int    `json:"timeoutSeconds"`
	PID            int    `json:"-"`
	ProcessGroup   int    `json:"-"`
}

type Result struct {
	OK    bool
	Error string
}

func Check(cfg Config) Result {
	timeout := time.Duration(cfg.TimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	deadline := time.Now().Add(timeout)

	switch cfg.Type {
	case "none":
		return Result{OK: true}

	case "process":
		return checkProcess(cfg)

	case "tcp":
		return retryUntil(deadline, func() Result {
			return checkTCP(cfg)
		})

	case "http":
		return retryUntil(deadline, func() Result {
			return checkHTTP(cfg)
		})

	default:
		return Result{OK: false, Error: fmt.Sprintf("unknown health check type: %s", cfg.Type)}
	}
}

func retryUntil(deadline time.Time, check func() Result) Result {
	for {
		result := check()
		if result.OK {
			return result
		}
		if time.Now().After(deadline) {
			return Result{
				OK:    false,
				Error: fmt.Sprintf("health check timed out: %s", result.Error),
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func checkProcess(cfg Config) Result {
	if cfg.PID <= 0 {
		return Result{OK: false, Error: "no PID"}
	}

	// Wait for grace period to see if process stays alive
	time.Sleep(500 * time.Millisecond)

	// Try to find process; if it exited, health check fails
	proc, err := os.FindProcess(cfg.PID)
	if err != nil {
		return Result{OK: false, Error: fmt.Sprintf("process %d not found", cfg.PID)}
	}

	// Signal 0 checks existence on Unix; on Windows FindProcess always succeeds
	if proc.Signal(syscall.Signal(0x0)) != nil {
		return Result{OK: false, Error: fmt.Sprintf("process %d exited", cfg.PID)}
	}

	return Result{OK: true}
}

func checkTCP(cfg Config) Result {
	host := cfg.Host
	if host == "" {
		host = "127.0.0.1"
	}
	addr := net.JoinHostPort(host, cfg.Port)

	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		return Result{OK: false, Error: err.Error()}
	}
	conn.Close()
	return Result{OK: true}
}

func checkHTTP(cfg Config) Result {
	client := &http.Client{Timeout: 2 * time.Second}

	req, err := http.NewRequest("GET", cfg.URL, nil)
	if err != nil {
		return Result{OK: false, Error: err.Error()}
	}

	resp, err := client.Do(req)
	if err != nil {
		return Result{OK: false, Error: err.Error()}
	}
	defer resp.Body.Close()

	if cfg.ExpectedStatus > 0 && resp.StatusCode != cfg.ExpectedStatus {
		return Result{OK: false, Error: fmt.Sprintf("expected status %d, got %d", cfg.ExpectedStatus, resp.StatusCode)}
	}

	return Result{OK: true}
}
