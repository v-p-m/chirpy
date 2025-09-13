package main

import (
	"log"
	"net/http"
)



func main() {
	const filepathRoot = "."
	const port = "8080"

	// Create a new request multiplexer (router) to handle HTTP requests.
	// This allows you to register multiple routes and handlers.
	mux := http.NewServeMux()

	// Create a file server to serve static files from the current directory (".").
	// This will serve files like index.html, CSS, JS, etc.
	filesrv := http.FileServer(http.Dir(filepathRoot))

	// Register the file server as the handler for the root path ("/").
	// This means all requests to the root and its subpaths will be handled by the file server.
	mux.Handle("/app/", http.StripPrefix("/app", filesrv))

	// Register the handler function for the "/healthz" endpoint.
	// This endpoint is typically used for health checks.
	mux.HandleFunc("/healthz", handlerHealthz)

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