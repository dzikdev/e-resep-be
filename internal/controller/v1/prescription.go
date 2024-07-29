package v1

import (
	"context"
	"e-resep-be/internal/config"
	"e-resep-be/internal/helper"
	"e-resep-be/internal/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	// PrescriptionController is an interface that has all the function to be implemented inside prescription controller
	PrescriptionController interface {
		Create(ctx echo.Context) error
	}

	// PrescriptionControllerImpl is an app prescription struct that consists of all the dependencies needed for prescription controller
	PrescriptionControllerImpl struct {
		Context context.Context
		Config  *config.Configuration
	}
)

// NewPrescriptionController return new instance prescription controller
func NewPrescriptionController(ctx context.Context, config *config.Configuration) *PrescriptionControllerImpl {
	return &PrescriptionControllerImpl{
		Context: ctx,
		Config:  config,
	}
}

func (pc *PrescriptionControllerImpl) Create(ctx echo.Context) error {
	var prescriptionReq model.PrescriptionRequest

	if err := ctx.Bind(&prescriptionReq); err != nil {
		return helper.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), err.Error(), err, nil)
	}

	//TODO : completing this logic

	return helper.NewResponses[any](ctx, http.StatusCreated, "Success Create Prescription", nil, nil, nil)
}
