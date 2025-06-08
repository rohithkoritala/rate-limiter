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
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	// Create limiter
	ratelimiter := limiter.NewInMemoryLimiter()

	// Set up API routes
	handler := api.NewHandler(ratelimiter)

	log.Println("Starting server on :8080...")
	http.ListenAndServe(":8080", handler)
}
