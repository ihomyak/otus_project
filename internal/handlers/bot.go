package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ihomyak/otus_project/internal/entities"
	"github.com/ihomyak/otus_project/internal/http/server"
	"github.com/jmoiron/sqlx"
)

func CreateBot(w http.ResponseWriter, r *http.Request, psql *sqlx.DB) {
	if r.Method != http.MethodPost {
		server.RespondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var bot entities.Bot
	err := json.NewDecoder(r.Body).Decode(&bot)
	if err != nil {
		server.RespondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if bot.ID == "" {
		server.RespondError(w, "ID is required", http.StatusBadRequest)
		return
	}
	if bot.Name == "" {
		server.RespondError(w, "Name is required", http.StatusBadRequest)
		return
	}

	bot.CreatedAt = time.Now()
	bot.UpdatedAt = time.Now()

	query := `INSERT INTO bots (id, name, created_at, updated_at) VALUES ($1, $2, $3, $4)`
	_, err = psql.Exec(query, bot.ID, bot.Name, bot.CreatedAt, bot.UpdatedAt)
	if err != nil {
		server.RespondError(w, "Error inserting bot into database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	server.RespondStructCreated(w, bot)
}

func UpdateBot(w http.ResponseWriter, r *http.Request, psql *sqlx.DB, id string) {
	var bot entities.Bot
	err := json.NewDecoder(r.Body).Decode(&bot)
	if err != nil {
		server.RespondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if bot.Name == "" {
		server.RespondError(w, "Name is required", http.StatusBadRequest)
		return
	}
	bot.ID = id
	bot.UpdatedAt = time.Now()
	query := `UPDATE bots SET name = $1, updated_at = $2 WHERE id = $3`
	_, err = psql.Exec(query, bot.Name, bot.UpdatedAt, id)
	if err != nil {
		server.RespondError(w, "Error updating bot in database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = psql.Get(&bot, "SELECT * FROM bots WHERE id = $1", id)
	if err != nil {
		server.RespondError(w, "Error with query: "+err.Error(), http.StatusInternalServerError)
		return
	}
	server.RespondStructOk(w, bot)
}

func GetBot(w http.ResponseWriter, psql *sqlx.DB, id string) {
	var bot entities.Bot
	err := psql.Get(&bot, "SELECT * FROM bots WHERE id = $1", id)
	if err != nil {
		server.RespondError(w, "Error with query: "+err.Error(), http.StatusInternalServerError)
		return
	}
	server.RespondStructOk(w, bot)
}

func DeleteBot(w http.ResponseWriter, psql *sqlx.DB, id string) {
	_, err := psql.Exec("DELETE FROM bots WHERE id = $1", id)
	if err != nil {
		server.RespondError(w, "Error deleting bot: "+err.Error(), http.StatusInternalServerError)
		return
	}
	server.RespondSuccess(w, "Bot deleted successfully")
}
