package entities

import (
	"database/sql"
	"time"
)

type Token struct {
	ID        string         `json:"id" db:"id"`
	BotID     string         `json:"bot_id" db:"bot_id"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
	Rate      int            `json:"rate" db:"rate"`
	Period    sql.NullString `json:"period" db:"period"`
}
