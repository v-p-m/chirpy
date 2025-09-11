package main

import (
	"log"
	"net/http"
)



func main() {
	// create a new http.ServerMux
	mux := http.NewServeMux()

	// create a new http.Server struct
	server := &http.Server {
		Addr:		":8080",	// Set the .Addr field to ":8080"
		Handler:	mux,		// Use the new "ServeMux" as the server's handler
	}

	// use the server's listenAndServer method to start the server
	log.Println("Server starting on port 8080...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}