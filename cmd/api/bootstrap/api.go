package bootstrap

import (
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

	// Registrando a rota
	app.Post("/api/health-check", healtchCheckHandler.HealthCheck)

	return nil
}
