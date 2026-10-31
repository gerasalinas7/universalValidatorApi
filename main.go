package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"universal-validator-api/internal/validator"
)

type validateRequest struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type validateResponse struct {
	Valid  bool     `json:"valid"`
	Type   string   `json:"type"`
	Value  string   `json:"value"`
	Errors []string `json:"errors"`
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("GET /api/v1/validate/types", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"types": []string{"email", "cuit", "phone", "url", "password"},
		})
	})

	mux.HandleFunc("POST /api/v1/validate", func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		defer r.Body.Close()

		var req validateRequest
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON request"})
			return
		}

		req.Type = validator.NormalizeType(req.Type)
		if req.Type == "" || req.Value == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Both 'type' and 'value' are required"})
			return
		}

		result, supported := validator.Validate(req.Type, req.Value)
		if !supported {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "Unsupported validation type. Check GET /api/v1/validate/types",
			})
			return
		}

		writeJSON(w, http.StatusOK, validateResponse{
			Valid: result.Valid, Type: req.Type, Value: req.Value, Errors: result.Errors,
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{Addr: ":" + port, Handler: withCORS(mux)}
	log.Printf("Universal Validator API listening on :%s", port)
	log.Fatal(server.ListenAndServe())
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
