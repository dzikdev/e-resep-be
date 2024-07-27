package application

import (
	controllerV1 "e-resep-be/internal/controller/v1"
	"e-resep-be/internal/repository"
	"e-resep-be/internal/service"
)

type Dependency struct {
	HealthCheckController controllerV1.HealthCheckController
}

func SetupDependencyInjection(app *App) *Dependency {
	// repository
	healthCheckRepoImpl := repository.NewHealthCheckRepository(app.Context, app.Config, app.Logger, app.DB)

	// service
	healthCheckSvcImpl := service.NewHealthCheckService(app.Context, app.Config, healthCheckRepoImpl)

	// controller
	healthCheckControllerImpl := controllerV1.NewHealthCheckController(app.Context, app.Config, healthCheckSvcImpl)

	return &Dependency{
		HealthCheckController: healthCheckControllerImpl,
	}
}
