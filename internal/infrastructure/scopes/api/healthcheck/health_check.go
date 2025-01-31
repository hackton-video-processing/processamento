package healthcheck

import (
	"net/http"
)

type (
	healthCheckUseCase interface {
		Check() string
	}

	HealthCheckHandler struct {
		healthUseCase healthCheckUseCase
	}
)

func NewHealthCheckHandler(healthCheckUseCase healthCheckUseCase) *HealthCheckHandler {
	return &HealthCheckHandler{
		healthUseCase: healthCheckUseCase,
	}
}

func (h *HealthCheckHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	answer := h.healthUseCase.Check()

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(answer))
}
