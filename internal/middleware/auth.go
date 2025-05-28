package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ihomyak/otus_project/internal/entities"
	"github.com/ihomyak/otus_project/internal/http/server"
	"github.com/ihomyak/otus_project/internal/services"
	"github.com/jmoiron/sqlx"
)

func AuthMiddleware(next http.HandlerFunc, psql *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			server.RespondError(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			server.RespondError(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		token := tokenParts[1]
		isValid, tokenEntity, err := isValidToken(token, psql)
		if err != nil || !isValid {
			server.RespondError(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// todo: add bot limit checks
		limiter := services.RateLimiter{
			TokenID:  tokenEntity.ID,
			Limit:    tokenEntity.Rate,
			Duration: services.ConvertTextToTime(tokenEntity.Period.String),
		}
		ctx := context.WithValue(r.Context(), struct{ name string }{services.RateLimitContextKey}, limiter)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func isValidToken(id string, psql *sqlx.DB) (bool, entities.Token, error) {
	var token entities.Token
	err := psql.Get(&token, "SELECT * FROM tokens WHERE id = $1 LIMIT 1", id)
	if err != nil {
		return false, token, err
	}

	return true, token, err
}
