package main

import (
	"fmt"
	"httpProject/handle"
	"log"
	"net/http"
)

const port = 9000

func main() {
	app := handle.NewApp()
	go app.SessionManager.SessionCleanup()
	// go app.SessionManager.SessionCleanup()

	log.Printf("Server starting on port: %d\n", port)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: app,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Could not start server: %v\n", err)
	}
}
