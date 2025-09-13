package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)
type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
    // Define the root directory for serving static files.
	const filepathRoot = "."
	// Define the port on which the server will listen.
	const port = "8080"
	// Initialize the API configuration struct.
	apiCfg := &apiConfig{}

	// Create a new request multiplexer (router) to handle HTTP requests.
	// This allows you to register multiple routes and handlers.
	mux := http.NewServeMux()

	// Create a file server to serve static files from the current directory (".").
	// This will serve files like index.html, CSS, JS, etc.
	filesrv := http.FileServer(http.Dir(filepathRoot))

    // Register the file server as the handler for the "/app/" path.
    // http.StripPrefix removes the "/app" prefix before passing the request to the file server.
    // The middlewareMetricsInc middleware increments the hit counter for each request.
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(filesrv)))

	// Register the handler function for the "/healthz" endpoint.
	// This endpoint is typically used for health checks.
	mux.HandleFunc("GET /api/healthz", handlerHealthz)

	// Register the handler function for the "/metrics" endpoint.
    // This endpoint returns the number of hits to the file server.
	mux.HandleFunc("GET /api/metrics", apiCfg.handlerMetrics)

	// Register the handler function for the "/reset" endpoint.
    // This endpoint resets the hit counter to zero.
	mux.HandleFunc("POST /api/reset", apiCfg.handlerMetricsRes)

	// Create a new HTTP server instance with custom settings.
	server := &http.Server {
		Addr:		":" + port,	// Set the server to listen on port 8080 for all network interfaces.
		Handler:	mux,		// Use the ServeMux as the server's request handler.
	}

	// Start the HTTP server and listen for incoming requests.
	// This is a blocking call, so the program will wait here until the server is stopped.
	log.Println("Server starting on port 8080...")
	if err := server.ListenAndServe(); err != nil {
		// If the server fails to start, log the error and exit.
		log.Fatalf("Server failed: %v", err)
	}
}

// handlerHealthz responds to health check requests with a 200 OK status.
// It sets the response content type to plain text and writes the status text.
func handlerHealthz(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// handlerMetrics responds to requests for the metrics endpoint.
// It returns the current number of hits to the file server.
func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hits := cfg.fileserverHits.Load()
	w.Write([]byte(fmt.Sprintf("Hits: %d", hits)))
}

// middlewareMetricsInc is a middleware that increments the hit counter for each request.
// It wraps the next handler and increments the counter before passing the request to it.
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cfg.fileserverHits.Add(1)
        next.ServeHTTP(w, r)
    })
}

// handlerMetricsRes responds to requests for the reset endpoint.
// It resets the hit counter to zero.
func (cfg *apiConfig) handlerMetricsRes(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}