package health

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckNone(t *testing.T) {
	r := Check(Config{Type: "none"})
	if !r.OK {
		t.Errorf("expected ok, got error: %s", r.Error)
	}
}

func TestCheckUnknown(t *testing.T) {
	r := Check(Config{Type: "invalid"})
	if r.OK {
		t.Error("expected failure for unknown type")
	}
}

func TestCheckHTTP(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	r := Check(Config{
		Type:           "http",
		URL:            server.URL,
		ExpectedStatus: 200,
		TimeoutSeconds: 5,
	})
	if !r.OK {
		t.Errorf("expected ok, got: %s", r.Error)
	}
}

func TestCheckHTTPWrongStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	r := Check(Config{
		Type:           "http",
		URL:            server.URL,
		ExpectedStatus: 200,
		TimeoutSeconds: 2,
	})
	if r.OK {
		t.Error("expected failure for wrong status")
	}
}

func TestCheckHTTPTimeout(t *testing.T) {
	r := Check(Config{
		Type:           "http",
		URL:            "http://192.0.2.1:1/",
		ExpectedStatus: 200,
		TimeoutSeconds: 1,
	})
	if r.OK {
		t.Error("expected timeout failure")
	}
}

func TestCheckTCPRefused(t *testing.T) {
	r := Check(Config{
		Type:           "tcp",
		Host:           "127.0.0.1",
		Port:           "1",
		TimeoutSeconds: 1,
	})
	if r.OK {
		t.Error("expected failure for refused connection")
	}
}

func TestStatusError(t *testing.T) {
	err := &StatusError{Expected: 200, Got: 404}
	msg := err.Error()
	if msg != "expected status 200, got 404" {
		t.Errorf("unexpected error message: %s", msg)
	}
}
