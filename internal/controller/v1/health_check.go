package v1

import (
	"context"
	"net/http"

	"e-resep-be/internal/config"
	"e-resep-be/internal/helper"
	"e-resep-be/internal/service"

	"github.com/labstack/echo/v4"
)

type (
	// HealthCheckController is an interface that has all the function to be implemented inside health check controller
	HealthCheckController interface {
		Check(ctx echo.Context) error
	}

	// HealthCheckControllerImpl is an app health check struct that consists of all the dependencies needed for health check controller
	HealthCheckControllerImpl struct {
		Context        context.Context
		Config         *config.Configuration
		HealthCheckSvc service.HealthCheckService
	}
)

// NewHealthCheckController return new instances health check controller
func NewHealthCheckController(ctx context.Context, config *config.Configuration, healthCheckSvc service.HealthCheckService) *HealthCheckControllerImpl {
	return &HealthCheckControllerImpl{
		Context:        ctx,
		Config:         config,
		HealthCheckSvc: healthCheckSvc,
	}
}

func (hc *HealthCheckControllerImpl) Check(ctx echo.Context) error {
	ok, err := hc.HealthCheckSvc.Check()
	if err != nil || !ok {
		return helper.NewResponses[any](ctx, http.StatusInternalServerError, "Not OK", ok, err, nil)
	}

	return helper.NewResponses[any](ctx, http.StatusOK, "OK", ok, nil, nil)
}
