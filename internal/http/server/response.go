package server

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

func RespondError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(Response{
		Ok:      false,
		Message: message,
	})
	if err != nil {
		log.Printf("response marshal error: %s", err)
	}
}

func RespondSuccess(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(Response{
		Ok:      true,
		Message: message,
	})
	if err != nil {
		log.Printf("response marshal error: %s", err)
	}
}

func RespondStructOk(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("response marshal error: %s", err)
	}
}

func RespondStructCreated(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("response marshal error: %s", err)
	}
}
