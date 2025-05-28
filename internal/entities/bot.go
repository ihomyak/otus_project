package entities

import (
	"database/sql"
	"time"
)

// Bot представляет модель бота.
type Bot struct {
	ID        string         `json:"id" db:"id"`
	Name      string         `json:"name" db:"name"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
	Rate      int            `json:"rate" db:"rate"`
	Period    sql.NullString `json:"period" db:"period"`
}
