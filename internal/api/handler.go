package api

import (
	"encoding/json"
	"log"
	"net/http"
	"rate-limiter/internal/limiter"
)

type Handler struct {
	limiter *limiter.Limiter
}

func NewHandler(l *limiter.Limiter) http.Handler {
	h := &Handler{limiter: l}
	mux := http.NewServeMux()
	mux.HandleFunc("/check", h.checkRateLimit)
	mux.HandleFunc("/config", h.setConfig)
	return mux
}

func (h *Handler) checkRateLimit(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "missing key", http.StatusBadRequest)
		return
	}

	if h.limiter.Allow(key) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("allowed"))
	} else {
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("rate limit exceeded"))
	}
}

type ConfigRequest struct {
	Key   string  `json:"key"`
	Rate  float64 `json:"rate"`
	Burst int     `json:"burst"`
}

func (h *Handler) setConfig(w http.ResponseWriter, r *http.Request) {
	var cfg ConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	h.limiter.SetRate(cfg.Key, cfg.Rate, cfg.Burst)
	log.Printf("Setting rate limit for %s: rate=%v, burst=%v", cfg.Key, cfg.Rate, cfg.Burst)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("rate limit config updated"))
}
