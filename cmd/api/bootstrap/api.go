package bootstrap

import (
	"github.com/hackton-video-processing/processamento/internal/infrastructure/scopes/api/createprocess"
	"github.com/hackton-video-processing/processamento/internal/infrastructure/scopes/api/getprocessbyid"
	"github.com/hackton-video-processing/processamento/internal/infrastructure/scopes/catalog"
	"log"

	"github.com/hackton-video-processing/processamento/internal/infrastructure/scopes/api/healthcheck"

	"github.com/go-chi/chi/v5"
)

func createAPIRoutes(app *chi.Mux, catalog catalog.UseCase) error {
	log.Println("Creating api routes")

	// Inicializando o handler
	healtchCheckHandler, err := healthcheck.BootStrapHealth(catalog)
	if err != nil {
		return err
	}

	createProcessHandler, err := createprocess.BootStrapCreateProcess(catalog)
	if err != nil {
		return err
	}

	getProcessByIDHandler, err := getprocessbyid.BootstrapGetProcessBtID(catalog)
	if err != nil {
		return err
	}

	// health check
	app.Get("/", healtchCheckHandler.HealthCheck)

	// endpoints
	app.Post("/api/process", createProcessHandler.CreateProcess)
	app.Get("/api/process/{id}", getProcessByIDHandler.GetProcessByID)

	return nil
}
