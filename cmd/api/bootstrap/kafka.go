package bootstrap

import (
	"log"

	"github.com/hackton-video-processing/processamento/internal/infrastructure/scopes/catalog"
	"github.com/hackton-video-processing/processamento/internal/infrastructure/scopes/kafka"

	"github.com/go-chi/chi/v5"
)

func createMSKConsumerRoutes(app *chi.Mux, catalog catalog.UseCase) error {
	log.Println("Creating kafka routes")

	// Inicializando o handler
	videoProcessingHandler, err := kafka.BootStrapVideoProcessing(catalog)
	if err != nil {
		return err
	}

	// Registrando a rota
	app.Post("/sink/process/video", videoProcessingHandler.VideoProcessing)

	return nil
}
