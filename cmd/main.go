package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/ihomyak/otus_project/config"
	"github.com/ihomyak/otus_project/db"
	"github.com/ihomyak/otus_project/internal/handlers"
	"github.com/ihomyak/otus_project/internal/http/server"
	"github.com/ihomyak/otus_project/internal/logger"
	"github.com/ihomyak/otus_project/internal/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
)

type HookRequest struct {
	Text string `json:"text"`
}

var defaultReadHeaderTimeout = 10 * time.Second

func main() {
	log.Println("Starting service...")
	appConfig, err := config.GetConfig()
	if err != nil {
		log.Fatal("Error setting up the configuration ", err)
	}

	l := logger.NewLogger(appConfig.Logger.Level)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	// Инициализация Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     appConfig.Redis.Addr,
		Password: appConfig.Redis.Password,
		DB:       appConfig.Redis.Database,
	})

	// Проверка подключения к Redis
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}

	psql, err := db.ConnectDB(appConfig)
	defer func(psql *sqlx.DB) {
		err := psql.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(psql)

	if err != nil {
		l.Error("Error to setup a Postgresql connection", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Service is up and running"))
	})

	mux.HandleFunc("/createPost", func(w http.ResponseWriter, r *http.Request) {
		middleware.AuthMiddleware(
			middleware.RateLimiterMiddleware(
				func(w http.ResponseWriter, r *http.Request) {
					handlers.CreatePostHandler(w, r)
				},
				rdb,
			),
			psql,
		).ServeHTTP(w, r)
	})

	mux.HandleFunc("/crud/bot", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateBot(w, r, psql)
	})

	mux.HandleFunc("/crud/bot/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		switch r.Method {
		case http.MethodDelete:
			handlers.DeleteBot(w, psql, id)
		case http.MethodPut:
			handlers.UpdateBot(w, r, psql, id)
		case http.MethodGet:
			handlers.GetBot(w, psql, id)
		default:
			server.RespondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/crud/token", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateToken(w, r, psql)
	})

	mux.HandleFunc("/crud/token/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		switch r.Method {
		case http.MethodGet:
			handlers.GetToken(w, psql, id)
		case http.MethodPut:
			handlers.UpdateToken(w, r, psql, id)
		case http.MethodDelete:
			handlers.DeleteToken(w, psql, id)
		default:
			server.RespondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/crud/hook", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateHook(w, r, psql)
	})

	mux.HandleFunc("/crud/hook/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		switch r.Method {
		case http.MethodGet:
			handlers.GetHook(w, psql, id)
		case http.MethodPut:
			handlers.UpdateHook(w, r, psql, id)
		case http.MethodDelete:
			handlers.DeleteHook(w, psql, id)
		default:
			server.RespondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.Handle("/metrics", promhttp.Handler())

	handler := middleware.MetricsMiddleware(mux)
	httpServer := &http.Server{
		Addr:              ":" + appConfig.Server.Port,
		Handler:           handler,
		ReadHeaderTimeout: defaultReadHeaderTimeout,
	}

	l.Info("server start on port", slog.String("port", appConfig.Server.Port))
	log.Panic(httpServer.ListenAndServe())
}
