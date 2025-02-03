package healthcheck

import (
	"fmt"
	"github.com/hackton-video-processing/processamento/internal/infrastructure/scopes/catalog"
	"github.com/hackton-video-processing/processamento/pkg/once"
)

func BootStrapHealth(useCaseCatalog catalog.UseCase) (*HealthCheckHandler, error) {
	healthCheckUsecase, err := once.Call(useCaseCatalog.Health)
	if err != nil {
		return nil, fmt.Errorf("creating healthcheck use case: %w", err)
	}

	return NewHealthCheckHandler(healthCheckUsecase), nil
}
