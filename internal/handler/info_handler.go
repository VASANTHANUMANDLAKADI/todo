package handler

import (
	"encoding/json"
	"net/http"
)

type InfoResponse struct {
	Service string `json:"service"`
	Version string `json:"version"`
	Status  string `json:"status"`
}

func InfoHandler(w http.ResponseWriter, r *http.Request) {

	response := InfoResponse{
		Service: "todo-api",
		Version: "1.0.0",
		Status:  "running",
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)

}
