package getprocessbyid

import (
	"fmt"
	"github.com/hackton-video-processing/processamento/internal/infrastructure/scopes/catalog"
	"github.com/hackton-video-processing/processamento/pkg/once"
)

func BootstrapGetProcessBtID(useCaseCatalog catalog.UseCase) (*GetProcessByIDHandler, error) {
	getProcessByIDUsecase, err := once.Call(useCaseCatalog.GetProcessByID)
	if err != nil {
		return nil, fmt.Errorf("creating process use case: %w", err)
	}

	return NewGetProcessByIDHandler(getProcessByIDUsecase), nil
}
