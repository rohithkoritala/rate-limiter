package main

import (
	"context"
	"log"
	"net/http"
	"rate-limiter/internal/api"
	"rate-limiter/internal/limiter"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    []string{"127.0.0.1:7000", "127.0.0.1:7001", "127.0.0.1:7002", "127.0.0.1:7003", "127.0.0.1:7004", "127.0.0.1:7005"},
		Password: "",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Could not connect to Redis Cluster: %v", err)
	}

	// Create Redis-based limiter
	ratelimiter := limiter.NewRedisClusterLimiter(rdb)
	ratelimiter.SetRate("user123", 1, 5)

	// Set up API routes
	handler := api.NewHandler(ratelimiter)

	log.Println("Starting server on :8080...")
	http.ListenAndServe(":8080", handler)
}
