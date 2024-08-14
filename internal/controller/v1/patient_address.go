package v1

import (
	"context"
	"e-resep-be/internal/config"
	"e-resep-be/internal/helper"
	"e-resep-be/internal/model"
	"e-resep-be/internal/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type (
	// PatientAddressController is an interface that has all the function to be implemented inside patient address controller
	PatientAddressController interface {
		Create(ctx echo.Context) error
		Update(ctx echo.Context) error
	}

	// PatientAddressControllerImpl is an app patient address struct that consists of all the dependencies needed for patient address controller
	PatientAddressControllerImpl struct {
		Context           context.Context
		Config            *config.Configuration
		PatientAddressSvc service.PatientAddressService
	}
)

// NewPrescriptionController return new instance patient address controller
func NewPatientAddressController(ctx context.Context, config *config.Configuration, patientAddressSvc service.PatientAddressService) *PatientAddressControllerImpl {
	return &PatientAddressControllerImpl{
		Context:           ctx,
		Config:            config,
		PatientAddressSvc: patientAddressSvc,
	}
}

func (pc *PatientAddressControllerImpl) Create(ctx echo.Context) error {
	var patientAddressReq model.CreateOrUpdatePatientAddressRequest

	if err := ctx.Bind(&patientAddressReq); err != nil {
		return helper.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), err.Error(), err, nil)
	}

	err := pc.PatientAddressSvc.Create(ctx.Request().Context(), &patientAddressReq)
	if err != nil {
		return helper.NewResponses[any](ctx, http.StatusInternalServerError, "Error Create Patient Address", nil, err, nil)
	}

	return helper.NewResponses[any](ctx, http.StatusCreated, "Success Create Patient Address", nil, nil, nil)
}

func (pc *PatientAddressControllerImpl) Update(ctx echo.Context) error {
	var patientAddressReq model.CreateOrUpdatePatientAddressRequest

	id := ctx.Param("id")

	parseID, err := strconv.Atoi(id)
	if err != nil {
		return helper.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), err.Error(), err, nil)
	}

	if err := ctx.Bind(&patientAddressReq); err != nil {
		return helper.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), err.Error(), err, nil)
	}

	err = pc.PatientAddressSvc.Update(ctx.Request().Context(), &patientAddressReq, parseID)
	if err != nil {
		return helper.NewResponses[any](ctx, http.StatusInternalServerError, "Error Update Patient Address", nil, err, nil)
	}

	return helper.NewResponses[any](ctx, http.StatusOK, "Success Update Patient Address", nil, nil, nil)
}
