package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ihomyak/otus_project/internal/http/server"
)

type PostRequest struct {
	Text string `json:"text"`
	Data string `json:"data"`
	IP   string `json:"ip"`
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		server.RespondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req PostRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		server.RespondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Text == "" {
		server.RespondError(w, "Text is required", http.StatusBadRequest)
		return
	}

	// Здесь можно добавить логику сохранения поста в базу данных

	server.RespondSuccess(w, "Post created successfully")
}
