package handler

import (
	"encoding/json"
	"net/http"
)

// apiResponse is the standard JSON envelope for every endpoint.
// On success:  { "data": <payload> }
// On failure:  { "error": "<message>" }
type apiResponse struct {
	Data  any    `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

// writeJSON serialises data into the success envelope and writes it to w.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(apiResponse{Data: data})
}

// writeError serialises a plain string into the error envelope and writes it to w.
func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(apiResponse{Error: msg})
}