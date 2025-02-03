package createprocess

import (
	"fmt"
	"github.com/hackton-video-processing/processamento/internal/infrastructure/scopes/catalog"
	"github.com/hackton-video-processing/processamento/pkg/once"
)

func BootStrapCreateProcess(useCaseCatalog catalog.UseCase) (*CreateProcessHandler, error) {
	createProcessUsecase, err := once.Call(useCaseCatalog.CreateProcess)
	if err != nil {
		return nil, fmt.Errorf("creating process use case: %w", err)
	}

	return NewCreateProcessHandler(createProcessUsecase), nil
}
