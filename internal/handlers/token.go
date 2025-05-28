package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ihomyak/otus_project/internal/entities"
	"github.com/ihomyak/otus_project/internal/http/server"
	"github.com/jmoiron/sqlx"
)

func CreateToken(w http.ResponseWriter, r *http.Request, psql *sqlx.DB) {
	if r.Method != http.MethodPost {
		server.RespondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var token entities.Token
	err := json.NewDecoder(r.Body).Decode(&token)
	if err != nil {
		server.RespondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if token.ID == "" {
		server.RespondError(w, "ID is required", http.StatusBadRequest)
		return
	}

	if token.BotID == "" {
		server.RespondError(w, "BotID is required", http.StatusBadRequest)
		return
	}

	token.CreatedAt = time.Now()
	token.UpdatedAt = time.Now()

	query := `INSERT INTO tokens (id, bot_id, created_at, updated_at) VALUES ($1, $2, $3, $4)`
	_, err = psql.Exec(query, token.ID, token.BotID, token.CreatedAt, token.UpdatedAt)
	if err != nil {
		server.RespondError(w, "Error inserting token into database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	server.RespondStructCreated(w, token)
}

func GetToken(w http.ResponseWriter, psql *sqlx.DB, id string) {
	var token entities.Token
	err := psql.Get(&token, "SELECT * FROM tokens WHERE id = $1", id)
	if err != nil {
		server.RespondError(w, "Error with query: "+err.Error(), http.StatusInternalServerError)
		return
	}
	server.RespondStructOk(w, token)
}

func UpdateToken(w http.ResponseWriter, r *http.Request, psql *sqlx.DB, id string) {
	var token entities.Token
	err := json.NewDecoder(r.Body).Decode(&token)
	if err != nil {
		server.RespondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if token.ID == "" {
		server.RespondError(w, "ID is required", http.StatusBadRequest)
		return
	}
	token.UpdatedAt = time.Now()
	query := `UPDATE tokens SET bot_id = $1, updated_at = $2 WHERE id = $3`
	_, err = psql.Exec(query, token.BotID, token.UpdatedAt, id)
	if err != nil {
		server.RespondError(w, "Error updating token in database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = psql.Get(&token, "SELECT * FROM tokens WHERE id = $1", id)
	if err != nil {
		server.RespondError(w, "Error with query: "+err.Error(), http.StatusInternalServerError)
		return
	}

	server.RespondStructOk(w, token)
}

func DeleteToken(w http.ResponseWriter, psql *sqlx.DB, id string) {
	_, err := psql.Exec("DELETE FROM tokens WHERE id = $1", id)
	if err != nil {
		server.RespondError(w, "Error deleting token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	server.RespondSuccess(w, "Token deleted successfully")
}
