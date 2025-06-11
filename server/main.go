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
		Addrs:    []string{"localhost:7000", "localhost:7001", "localhost:7002"},
		Password: "",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Could not connect to Redis Cluster: %v", err)
	}

	// Create Redis-based limiter
	ratelimiter := limiter.NewRedisClusterLimiter(rdb)

	// Set up API routes
	handler := api.NewHandler(ratelimiter)

	log.Println("Starting server on :8080...")
	http.ListenAndServe(":8080", handler)
}
