package bootstrap

import (
	"fmt"
	"github.com/hackton-video-processing/processamento/internal/infrastructure/aws/mysql"
	"github.com/hackton-video-processing/processamento/internal/infrastructure/aws/s3"
	"github.com/hackton-video-processing/processamento/internal/infrastructure/scopes/catalog"
	"log"
	"net/http"

	"github.com/hackton-video-processing/processamento/internal/infrastructure/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func CreateApplication() (*http.Server, error) {
	// Carregar configuração
	appConfig, err := config.LoadConfiguration()
	if err != nil {
		return nil, fmt.Errorf("error loading configuration: %w", err)
	}

	s3Client := s3.BootstrapS3(appConfig)

	db, err := mysql.BootstrapMySQLRepository(appConfig)
	if err != nil {
		return nil, fmt.Errorf("error bootstrapping MySQL repository: %w", err)
	}

	usecaseCatalog := catalog.New(appConfig, s3Client, db)

	// Criar roteador
	router := chi.NewRouter()

	// Configurar middlewares
	router.Use(middleware.Logger)            // Logger padrão
	router.Use(middleware.Recoverer)         // Recupera de panics
	router.Use(middleware.RequestID)         // Adiciona um ID único às requisições
	router.Use(middleware.Timeout(60 * 1e9)) // Define timeout de requisições

	// Configurar rotas
	err = SetupRoutes(router, usecaseCatalog)
	if err != nil {
		return nil, fmt.Errorf("error setting up routes: %w", err)
	}

	// Criar servidor HTTP
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", appConfig.Port),
		Handler: router,
	}

	log.Printf("Server starting on port %s in %s appConfig...", appConfig.Port, appConfig.Env)

	return server, nil
}
