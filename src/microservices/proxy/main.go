package main

import (

	"fmt"

	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var (
	monolithURL      string
	moviesServiceURL string
	eventsServiceURL string
	gradualMigration bool
	migrationPercent int
)

func main() {
	// Read environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	monolithURL = os.Getenv("MONOLITH_URL")
	if monolithURL == "" {
		monolithURL = "http://monolith:8080"
	}

	moviesServiceURL = os.Getenv("MOVIES_SERVICE_URL")
	if moviesServiceURL == "" {
		moviesServiceURL = "http://movies-service:8081"
	}

	eventsServiceURL = os.Getenv("EVENTS_SERVICE_URL")
	if eventsServiceURL == "" {
		eventsServiceURL = "http://events-service:8082"
	}

	gradualMigrationStr := os.Getenv("GRADUAL_MIGRATION")
	gradualMigration = gradualMigrationStr == "true"

	migrationPercentStr := os.Getenv("MOVIES_MIGRATION_PERCENT")
	if migrationPercentStr != "" {
		var err error
		migrationPercent, err = strconv.Atoi(migrationPercentStr)
		if err != nil {
			log.Printf("Invalid MOVIES_MIGRATION_PERCENT, defaulting to 0: %v", err)
			migrationPercent = 0
		}
	}

	// Custom handler that checks path before routing
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		log.Printf("Received request: %s %s", r.Method, path)
		
		// Handle specific routes first - order matters!
		// More specific routes must be checked before more general ones
		switch {
		case path == "/health":
			healthHandler(w, r)
			return
		case path == "/api/movies/health":
			log.Printf("Routing /api/movies/health to movies-service")
			handleMoviesHealth(w, r)
			return
		case path == "/api/movies":
			handleMovies(w, r)
			return
		case strings.HasPrefix(path, "/api/events/"):
			handleEventsProxy(w, r)
			return
		case strings.HasPrefix(path, "/api/movies/"):
			handleMovies(w, r)
			return
		case strings.HasPrefix(path, "/api/"):
			handleMonolithProxy(w, r)
			return
		default:
			http.NotFound(w, r)
			return
		}
	})

	log.Printf("Starting proxy service on port %s", port)
	log.Printf("Monolith URL: %s", monolithURL)
	log.Printf("Movies Service URL: %s", moviesServiceURL)
	log.Printf("Events Service URL: %s", eventsServiceURL)
	log.Printf("Gradual Migration: %v, Migration Percent: %d%%", gradualMigration, migrationPercent)

	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Strangler Fig Proxy is healthy"))
}

func handleMoviesHealth(w http.ResponseWriter, r *http.Request) {
	// Health check always goes to movies-service
	target, err := url.Parse(moviesServiceURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid movies service URL: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Routing /api/movies/health to movies-service")
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ServeHTTP(w, r)
}

func handleMovies(w http.ResponseWriter, r *http.Request) {
	// Strangler Fig pattern: route based on migration percentage
	var targetURL string
	var serviceName string

	if gradualMigration {
		// Random decision based on migration percentage
		randomValue := rand.Intn(100)
		if randomValue < migrationPercent {
			// Route to new movies service
			targetURL = moviesServiceURL
			serviceName = "movies-service"
		} else {
			// Route to monolith
			targetURL = monolithURL
			serviceName = "monolith"
		}
	} else {
		// If gradual migration is disabled, default to monolith
		targetURL = monolithURL
		serviceName = "monolith"
	}

	log.Printf("Routing /api/movies to %s (migration: %d%%)", serviceName, migrationPercent)

	// Create reverse proxy
	target, err := url.Parse(targetURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid target URL: %v", err), http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ServeHTTP(w, r)
}

func handleEventsProxy(w http.ResponseWriter, r *http.Request) {
	// Proxy all /api/events/* requests to events service
	target, err := url.Parse(eventsServiceURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid events service URL: %v", err), http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	// Modify the request path to remove /api prefix
	originalPath := r.URL.Path
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api")
	proxy.ServeHTTP(w, r)
	// Restore original path for logging
	r.URL.Path = originalPath
}

func handleMonolithProxy(w http.ResponseWriter, r *http.Request) {
	log.Printf("handleMonolithProxy called with path: %s", r.URL.Path)
	// Proxy all other /api/* requests to monolith
	target, err := url.Parse(monolithURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid monolith URL: %v", err), http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ServeHTTP(w, r)
}
