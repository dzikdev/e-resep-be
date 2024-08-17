package application

import (
	controllerV1 "e-resep-be/internal/controller/v1"
	"e-resep-be/internal/repository"
	"e-resep-be/internal/requester"
	"e-resep-be/internal/service"
)

type Dependency struct {
	HealthCheckController    controllerV1.HealthCheckController
	PrescriptionController   controllerV1.PrescriptionController
	AddressController        controllerV1.AddressController
	PatientAddressController controllerV1.PatientAddressController
	PaymentController        controllerV1.PaymentController
}

func SetupDependencyInjection(app *App) *Dependency {
	// requester
	whatsappRequesterImpl := requester.NewWhatsappRequester(app.Context, app.Config, app.Logger, app.HTTPClient)
	kimiaFarmaRequesterImpl := requester.NewKimiaFarmaRequester(app.Context, app.Config, app.Logger, app.HTTPClient)
	xenditRequesterImpl := requester.NewXenditRequester(app.Context, app.Config, app.Logger, app.XenditSDK)

	// repository
	healthCheckRepoImpl := repository.NewHealthCheckRepository(app.Context, app.Config, app.Logger, app.DB)
	prescriptionRepoImpl := repository.NewPrescriptionRepository(app.Context, app.Config, app.Logger, app.DB)
	addressRepoImpl := repository.NewAddressRepository(app.Context, app.Config, app.Logger, app.DB)
	patientRepoImpl := repository.NewPatientRepository(app.Context, app.Config, app.Logger, app.DB)
	patientAddressRepoImpl := repository.NewPatientAddressRepository(app.Context, app.Config, app.Logger, app.DB)
	medicationRepoImpl := repository.NewMedicationRepository(app.Context, app.Config, app.Logger, app.DB)
	transactionRepoImpl := repository.NewTransactionRepository(app.Context, app.Config, app.Logger, app.DB)
	paymentRepoImpl := repository.NewPaymentRepository(app.Context, app.Config, app.Logger, app.DB)

	// service
	healthCheckSvcImpl := service.NewHealthCheckService(app.Context, app.Config, healthCheckRepoImpl)
	prescriptionSvcImpl := service.NewPrescriptionService(app.Context, app.Config, prescriptionRepoImpl, whatsappRequesterImpl, kimiaFarmaRequesterImpl)
	addressSvcImpl := service.NewAddressService(app.Context, app.Config, addressRepoImpl)
	patientAddressSvcImpl := service.NewPatientAddressService(app.Context, app.Config, patientRepoImpl, patientAddressRepoImpl)
	paymentSvc := service.NewPaymentService(app.Context, app.Config, medicationRepoImpl, patientRepoImpl, patientAddressRepoImpl, transactionRepoImpl, paymentRepoImpl, kimiaFarmaRequesterImpl, xenditRequesterImpl)

	// controller
	healthCheckControllerImpl := controllerV1.NewHealthCheckController(app.Context, app.Config, healthCheckSvcImpl)
	prescriptionControllerImpl := controllerV1.NewPrescriptionController(app.Context, app.Config, prescriptionSvcImpl)
	addressControllerImpl := controllerV1.NewAddressController(app.Context, app.Config, addressSvcImpl)
	patientAddressControllerImpl := controllerV1.NewPatientAddressController(app.Context, app.Config, patientAddressSvcImpl)
	paymentControllerImpl := controllerV1.NewPaymentController(app.Context, app.Config, paymentSvc)

	return &Dependency{
		HealthCheckController:    healthCheckControllerImpl,
		PrescriptionController:   prescriptionControllerImpl,
		AddressController:        addressControllerImpl,
		PatientAddressController: patientAddressControllerImpl,
		PaymentController:        paymentControllerImpl,
	}
}
