package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ihomyak/otus_project/internal/entities"
	"github.com/ihomyak/otus_project/internal/http/server"
	"github.com/jmoiron/sqlx"
)

func CreateHook(w http.ResponseWriter, r *http.Request, psql *sqlx.DB) {
	if r.Method != http.MethodPost {
		server.RespondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var hook entities.Hook
	err := json.NewDecoder(r.Body).Decode(&hook)
	if err != nil {
		server.RespondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if hook.Name == "" {
		server.RespondError(w, "Name is required", http.StatusBadRequest)
		return
	}

	hook.CreatedAt = time.Now()
	hook.UpdatedAt = time.Now()

	query := `INSERT INTO hooks (id, name, created_at, updated_at) VALUES ($1, $2, $3, $4)`
	_, err = psql.Exec(query, hook.ID, hook.Name, hook.CreatedAt, hook.UpdatedAt)
	if err != nil {
		server.RespondError(w, "Error inserting hook into database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	server.RespondStructCreated(w, hook)
}

func GetHook(w http.ResponseWriter, psql *sqlx.DB, id string) {
	var hook entities.Hook
	err := psql.Get(&hook, "SELECT * FROM hooks WHERE id = $1", id)
	if err != nil {
		server.RespondError(w, "Error with query: "+err.Error(), http.StatusInternalServerError)
		return
	}
	server.RespondStructOk(w, hook)
}

func UpdateHook(w http.ResponseWriter, r *http.Request, psql *sqlx.DB, id string) {
	var hook entities.Hook
	err := json.NewDecoder(r.Body).Decode(&hook)
	if err != nil {
		server.RespondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if hook.Name == "" {
		server.RespondError(w, "Name is required", http.StatusBadRequest)
		return
	}
	hook.UpdatedAt = time.Now()
	query := `UPDATE hooks SET name = $1, updated_at = $2 WHERE id = $3`
	_, err = psql.Exec(query, hook.Name, hook.UpdatedAt, id)
	if err != nil {
		server.RespondError(w, "Error updating hook in database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	server.RespondStructOk(w, hook)
}

func DeleteHook(w http.ResponseWriter, psql *sqlx.DB, id string) {
	_, err := psql.Exec("DELETE FROM hooks WHERE id = $1", id)
	if err != nil {
		server.RespondError(w, "Error deleting hook: "+err.Error(), http.StatusInternalServerError)
		return
	}
	server.RespondSuccess(w, "Hook deleted successfully")
}
