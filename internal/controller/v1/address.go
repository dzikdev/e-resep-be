package v1

import (
	"context"
	"e-resep-be/internal/config"
	"e-resep-be/internal/helper"
	"e-resep-be/internal/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type (
	// AddressController is an interface that has all the function to be implemented inside address controller
	AddressController interface {
		GetProvince(ctx echo.Context) error
		GetCityByProvinceID(ctx echo.Context) error
		GetDistrictByCityID(ctx echo.Context) error
		GetSubDistrictByDistrictID(ctx echo.Context) error
	}

	// AddressControllerImpl is an app address struct that consists of all the dependencies needed for address controller
	AddressControllerImpl struct {
		Context    context.Context
		Config     *config.Configuration
		AddressSvc service.AddressService
	}
)

// NewAddressController return new instances address controller
func NewAddressController(ctx context.Context, config *config.Configuration, addressSvc service.AddressService) *AddressControllerImpl {
	return &AddressControllerImpl{
		Context:    ctx,
		Config:     config,
		AddressSvc: addressSvc,
	}
}

func (ac *AddressControllerImpl) GetProvince(ctx echo.Context) error {
	results, err := ac.AddressSvc.GetProvince(ctx.Request().Context())
	if err != nil {
		return helper.NewResponses[any](ctx, http.StatusInternalServerError, "Error Get Province", nil, err, nil)
	}

	return helper.NewResponses[any](ctx, http.StatusOK, "Success Get Province", results, nil, nil)
}

func (ac *AddressControllerImpl) GetCityByProvinceID(ctx echo.Context) error {
	provinceID := ctx.Param("id")

	parseProvinceID, err := strconv.Atoi(provinceID)
	if err != nil {
		return helper.NewResponses[any](ctx, http.StatusBadRequest, "Invalid ID", nil, err, nil)
	}

	results, err := ac.AddressSvc.GetCityByProvinceID(ctx.Request().Context(), parseProvinceID)
	if err != nil {
		return helper.NewResponses[any](ctx, http.StatusInternalServerError, "Error Get City", nil, err, nil)
	}

	return helper.NewResponses[any](ctx, http.StatusOK, "Success Get City", results, nil, nil)
}

func (ac *AddressControllerImpl) GetDistrictByCityID(ctx echo.Context) error {
	cityID := ctx.Param("id")

	parseCityID, err := strconv.Atoi(cityID)
	if err != nil {
		return helper.NewResponses[any](ctx, http.StatusBadRequest, "Invalid ID", nil, err, nil)
	}

	results, err := ac.AddressSvc.GetDistrictByCityID(ctx.Request().Context(), parseCityID)
	if err != nil {
		return helper.NewResponses[any](ctx, http.StatusInternalServerError, "Error Get District", nil, err, nil)
	}

	return helper.NewResponses[any](ctx, http.StatusOK, "Success Get District", results, nil, nil)
}

func (ac *AddressControllerImpl) GetSubDistrictByDistrictID(ctx echo.Context) error {
	districtID := ctx.Param("id")

	parseDistrictID, err := strconv.Atoi(districtID)
	if err != nil {
		return helper.NewResponses[any](ctx, http.StatusBadRequest, "Invalid ID", nil, err, nil)
	}

	results, err := ac.AddressSvc.GetSubDistrictByDistrictID(ctx.Request().Context(), parseDistrictID)
	if err != nil {
		return helper.NewResponses[any](ctx, http.StatusInternalServerError, "Error Get Sub District", nil, err, nil)
	}

	return helper.NewResponses[any](ctx, http.StatusOK, "Success Get Sub District", results, nil, nil)
}
