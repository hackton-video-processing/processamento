package bootstrap

import (
	"github.com/go-chi/chi/v5"
	"github.com/hackton-video-processing/processamento/internal/infrastructure/scopes/catalog"
)

func SetupRoutes(router *chi.Mux, catalog catalog.UseCase) error {
	if err := createAPIRoutes(router, catalog); err != nil {
		return err
	}

	if err := createMSKConsumerRoutes(router, catalog); err != nil {
		return err
	}

	return nil
}
