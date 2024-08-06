package application

import (
	controllerV1 "e-resep-be/internal/controller/v1"
	"e-resep-be/internal/repository"
	"e-resep-be/internal/service"
)

type Dependency struct {
	HealthCheckController  controllerV1.HealthCheckController
	PrescriptionController controllerV1.PrescriptionController
}

func SetupDependencyInjection(app *App) *Dependency {
	// requester

	// repository
	healthCheckRepoImpl := repository.NewHealthCheckRepository(app.Context, app.Config, app.Logger, app.DB)
	prescriptionRepoImpl := repository.NewPrescriptionRepository(app.Context, app.Config, app.Logger, app.DB)

	// service
	healthCheckSvcImpl := service.NewHealthCheckService(app.Context, app.Config, healthCheckRepoImpl)
	prescriptionSvcImpl := service.NewPrescriptionService(app.Context, app.Config, prescriptionRepoImpl)

	// controller
	healthCheckControllerImpl := controllerV1.NewHealthCheckController(app.Context, app.Config, healthCheckSvcImpl)
	prescriptionControllerImpl := controllerV1.NewPrescriptionController(app.Context, app.Config, prescriptionSvcImpl)

	return &Dependency{
		HealthCheckController:  healthCheckControllerImpl,
		PrescriptionController: prescriptionControllerImpl,
	}
}
