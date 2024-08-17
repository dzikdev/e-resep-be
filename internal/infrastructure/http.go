package infrastructure

import (
	"e-resep-be/internal/application"

	"github.com/labstack/echo/v4"
)

// ServeHTTP is wrapper function to start the apps infra in HTTP mode
func ServeHTTP(app *application.App) *echo.Echo {
	// call setup router
	setupRouter(app)

	return app.Application
}

// setupRouter is function to manage all routings
func setupRouter(app *application.App) {
	var dep = application.SetupDependencyInjection(app)

	v1 := app.Application.Group("/api/v1")
	{
		v1.GET("/health-check", dep.HealthCheckController.Check)

		prescription := v1.Group("/prescription")
		{
			prescription.POST("", dep.PrescriptionController.Create)

			prescription.GET("/:id", dep.PrescriptionController.GetByPrescriptionID)
		}

		v1.GET("/province", dep.AddressController.GetProvince)
		v1.GET("/province/:id/cities", dep.AddressController.GetCityByProvinceID)
		v1.GET("/cities/:id/district", dep.AddressController.GetDistrictByCityID)
		v1.GET("/district/:id/sub-district", dep.AddressController.GetSubDistrictByDistrictID)

		patient := v1.Group("/patient")
		{
			patient.POST("/address", dep.PatientAddressController.Create)
			patient.PUT("/address/:id", dep.PatientAddressController.Update)
		}

		payment := v1.Group("/payment")
		{
			payment.POST("/info", dep.PaymentController.GeneratePaymentInfo)
			payment.POST("", dep.PaymentController.CreatePayment)
			payment.POST("/notification", dep.PaymentController.PaymentNotification)
		}

	}
}
