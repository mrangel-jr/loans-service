package main

import (
	"log"
	"net/http"
	"os"

	"github.com/mrangel-jr/loans-service/controllers"
)

func main() {
	// Set up the logger
	logger := log.New(os.Stdout, "loansservice: ", log.LstdFlags)

	// Create a new HTTP server
	mux := http.NewServeMux()

	controllers.SetupRoutes(mux)

	// Start the server
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	logger.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("ListenAndServe: %v", err)
	}
}
