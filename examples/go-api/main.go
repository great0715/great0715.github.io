// Release Lab API is an illustrative, self-contained Go HTTP service.
// All data is held in memory. It is not used by the static portfolio demo.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

type release struct {
	ID          int64     `json:"id"`
	Service     string    `json:"service"`
	Environment string    `json:"environment"`
	Version     string    `json:"version"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type createRelease struct {
	Service     string `json:"service"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

type store struct {
	mu       sync.RWMutex
	nextID   int64
	releases []release
}

func newStore() *store { return &store{nextID: 1} }

func (s *store) list() []release {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]release, len(s.releases))
	copy(out, s.releases)
	return out
}

func (s *store) add(input createRelease) release {
	s.mu.Lock()
	defer s.mu.Unlock()
	item := release{
		ID:          s.nextID,
		Service:     input.Service,
		Environment: input.Environment,
		Version:     input.Version,
		Status:      "queued",
		CreatedAt:   time.Now().UTC(),
	}
	s.nextID++
	s.releases = append(s.releases, item)
	return item
}

func sendJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("encode response", "error", err)
	}
}

func readInput(r *http.Request) (createRelease, error) {
	var input createRelease
	r.Body = http.MaxBytesReader(nil, r.Body, 4096)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		return input, fmt.Errorf("invalid JSON: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return input, errors.New("expected exactly one JSON object")
	}
	input.Service = strings.TrimSpace(input.Service)
	input.Environment = strings.TrimSpace(input.Environment)
	input.Version = strings.TrimSpace(input.Version)
	if input.Service == "" || len(input.Service) > 80 || input.Version == "" || len(input.Version) > 80 {
		return input, errors.New("service and version are required (maximum 80 characters each)")
	}
	if input.Environment != "staging" && input.Environment != "production" {
		return input, errors.New("environment must be staging or production")
	}
	return input, nil
}

func handler(s *store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		sendJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /api/releases", func(w http.ResponseWriter, r *http.Request) {
		sendJSON(w, http.StatusOK, map[string]any{"releases": s.list()})
	})
	mux.HandleFunc("POST /api/releases", func(w http.ResponseWriter, r *http.Request) {
		if contentType := r.Header.Get("Content-Type"); !strings.HasPrefix(contentType, "application/json") {
			sendJSON(w, http.StatusUnsupportedMediaType, map[string]string{"error": "Content-Type must be application/json"})
			return
		}
		input, err := readInput(r)
		if err != nil {
			sendJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		sendJSON(w, http.StatusCreated, s.add(input))
	})
	return mux
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           handler(newStore()),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Error("shutdown", "error", err)
		}
	}()
	slog.Info("api listening", "port", port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("http server failed", "error", err)
		os.Exit(1)
	}
}
