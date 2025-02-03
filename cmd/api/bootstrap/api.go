package bootstrap

import (
	"github.com/hackton-video-processing/processamento/internal/infrastructure/scopes/api/createprocess"
	"github.com/hackton-video-processing/processamento/internal/infrastructure/scopes/api/getprocessbyid"
	"log"

	"github.com/hackton-video-processing/processamento/internal/infrastructure/config"
	"github.com/hackton-video-processing/processamento/internal/infrastructure/scopes/api/healthcheck"

	"github.com/go-chi/chi/v5"
)

func createAPIRoutes(app *chi.Mux, appConfig config.AppConfig) error {
	log.Println("Creating api routes")

	// Inicializando o handler
	healtchCheckHandler, err := healthcheck.BootStrapHealth(appConfig)
	if err != nil {
		return err
	}

	createProcessHandler, err := createprocess.BootStrapCreateProcess(appConfig)
	if err != nil {
		return err
	}

	getProcessByIDHandler, err := getprocessbyid.BootstrapGetProcessBtID(appConfig)
	if err != nil {
		return err
	}

	// health check
	app.Post("/api/health-check", healtchCheckHandler.HealthCheck)

	// endpoints
	app.Post("/api/process", createProcessHandler.CreateProcess)
	app.Get("/api/process/{id}", getProcessByIDHandler.GetProcessByID)

	return nil
}
