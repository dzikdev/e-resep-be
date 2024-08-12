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

			prescription.GET("/:id", dep.PrescriptionController.GetByID)
		}

		address := v1.Group("/address")
		{
			address.GET("/province", dep.AddressController.GetProvince)
			address.GET("/province/:id/cities", dep.AddressController.GetCityByProvinceID)
			address.GET("/cities/:id/district", dep.AddressController.GetDistrictByCityID)
			address.GET("/district/:id/sub-district", dep.AddressController.GetSubDistrictByDistrictID)
		}
	}
}
