package main

import (
	"log"
	"net/http"
)



func main() {
	// Create a new request multiplexer (router) to handle HTTP requests.
	// This allows you to register multiple routes and handlers.
	mux := http.NewServeMux()

	// Create a file server to serve static files from the current directory (".").
	// This will serve files like index.html, CSS, JS, etc.
	filesrv := http.FileServer(http.Dir("."))

	// Register the file server as the handler for the root path ("/").
	// This means all requests to the root and its subpaths will be handled by the file server.
	mux.Handle("/", filesrv)

	// Create a new HTTP server instance with custom settings.
	server := &http.Server {
		Addr:		":8080",	// Set the server to listen on port 8080 for all network interfaces.
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