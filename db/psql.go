package db

import (
	"github.com/ihomyak/otus_project/config"
	"github.com/jmoiron/sqlx" //nolint:nolintlint
	_ "github.com/lib/pq"     //nolint:gci
)

func ConnectDB(c *config.AppConfig) (*sqlx.DB, error) {
	return sqlx.Connect("postgres", c.Database.DSN)
}
