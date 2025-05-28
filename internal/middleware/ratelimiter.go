package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ihomyak/otus_project/internal/services"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
)

var (
	ctx           = context.Background()
	requestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "rate_limiter_blocks_total",
		Help: "Number of times requests were blocked by rate limiter",
	}, []string{"url", "uuid", "type"},
	)
)

func RateLimiterMiddleware(next http.HandlerFunc, rdb *redis.Client) http.HandlerFunc {
	script := `
	local key = KEYS[1]
	local now = tonumber(ARGV[1])
	local window = tonumber(ARGV[2])
	local limit = tonumber(ARGV[3])
	
	-- Удаляем записи вне временного окна
	redis.call('ZREMRANGEBYSCORE', key, 0, now - window)
	
	-- Получаем текущее количество запросов
	local current = redis.call('ZCARD', key)
	
	if current >= limit then
		return {0, current}
	end
	
	-- Добавляем новый запрос
	redis.call('ZADD', key, now, now)
	redis.call('EXPIRE', key, window/1000)
	
	return {1, current + 1}
	`

	sha, err := rdb.ScriptLoad(ctx, script).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to load rate limiter script: %v", err))
	}

	if err := prometheus.Register(requestsTotal); err != nil {
		panic(fmt.Sprintf("Failed to load prometheus metrics: %v", err))
	}

	return func(w http.ResponseWriter, r *http.Request) {
		rate, ok := r.Context().Value(services.RateLimitContextKey).(services.RateLimiter)
		if !ok {
			http.Error(w, "User not found in context", http.StatusInternalServerError)
			return
		}

		// todo: id_token, IP ???
		key := fmt.Sprintf("rate_limit:%s:%s", rate.TokenID, "192.168.0.1")
		limit := rate.Limit
		if limit <= 0 {
			limit = -1
		}

		now := time.Now().UnixMilli()
		window := rate.Duration

		result, err := rdb.EvalSha(ctx, sha, []string{key}, now, window.Milliseconds(), limit).Result()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		res := result.([]interface{})
		allowed := res[0].(int64) == 1
		current := res[1].(int64)

		if !allowed {
			requestsTotal.WithLabelValues(r.URL.Path, rate.TokenID, "block").Inc()
			w.Header().Set("Retry-After", "60")
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
			w.Header().Set("X-RateLimit-Remaining", "0")
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		requestsTotal.WithLabelValues(r.URL.Path, rate.TokenID, "access").Inc()
		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
		w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(limit-int(current)))

		next.ServeHTTP(w, r)
	}
}
