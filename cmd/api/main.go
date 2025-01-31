package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/hackton-video-processing/processamento/cmd/api/bootstrap"
)

func main() {
	if err := run(); err != nil {
		log.Printf("error running app: %v", err)
	}
}

func run() error {
	app, err := bootstrap.CreateApplication()
	if err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}

	// Iniciar servidor
	if err := app.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server error: %v", err)
	}

	return nil
}
