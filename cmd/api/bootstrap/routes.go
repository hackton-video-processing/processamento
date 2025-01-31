package bootstrap

import (
	"github.com/go-chi/chi/v5"
	"github.com/hackton-video-processing/processamento/internal/infrastructure/config"
)

func SetupRoutes(router *chi.Mux, appConfig config.AppConfig) error {
	if err := createAPIRoutes(router, appConfig); err != nil {
		return err
	}

	if err := createMSKConsumerRoutes(router, appConfig); err != nil {
		return err
	}

	return nil
}
