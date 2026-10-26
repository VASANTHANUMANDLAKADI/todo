package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInfoHandler(t *testing.T) {

	req := httptest.NewRequest(
		http.MethodGet,
		"/info",
		nil,
	)

	rec := httptest.NewRecorder()

	InfoHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, but got %d",
			http.StatusOK,
			rec.Code,
		)
	}

	if rec.Header().Get("Content-Type") != "application/json" {
		t.Errorf(
			"expected Content-Type application/json, but got %s",
			rec.Header().Get("Content-Type"),
		)
	}

	var response InfoResponse

	err := json.NewDecoder(rec.Body).Decode(&response)

	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Service != "todo-api" {
		t.Errorf("expected service todo-api, got %s",
			response.Service)
	}
	if response.Version != "1.0.0" {
		t.Errorf("expected version 1.0.0, got %s",
			response.Version)
	}
	if response.Status != "running" {
		t.Errorf("expected status running, got %s",
			response.Status)
	}
}
