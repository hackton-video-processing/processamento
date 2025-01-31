package bootstrap

import (
	"log"

	"github.com/hackton-video-processing/processamento/internal/infrastructure/config"
	"github.com/hackton-video-processing/processamento/internal/infrastructure/scopes/msk"

	"github.com/go-chi/chi/v5"
)

func createMSKConsumerRoutes(app *chi.Mux, appConfig config.AppConfig) error {
	log.Println("Creating msk routes")

	// Inicializando o handler
	videoProcessingHandler, err := msk.BootStrapVideoProcessing(appConfig)
	if err != nil {
		return err
	}

	// Registrando a rota
	app.Post("/consumer/process/video", videoProcessingHandler.VideoProcessing)

	return nil
}
