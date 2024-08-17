package v1

import (
	"context"
	"e-resep-be/internal/config"
	"e-resep-be/internal/helper"
	"e-resep-be/internal/model"
	"e-resep-be/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	// PaymentController is an interface that has all the function to be implemented inside payment controller
	PaymentController interface {
		GeneratePaymentInfo(ctx echo.Context) error
		CreatePayment(ctx echo.Context) error
		PaymentNotification(ctx echo.Context) error
	}

	// PaymentControllerImpl is an app payment struct that consists of all the dependencies needed for payment controller
	PaymentControllerImpl struct {
		Context    context.Context
		Config     *config.Configuration
		PaymentSvc service.PaymentService
	}
)

// NewPaymentController return new instance payment controller
func NewPaymentController(ctx context.Context, config *config.Configuration, paymentSvc service.PaymentService) *PaymentControllerImpl {
	return &PaymentControllerImpl{
		Context:    ctx,
		Config:     config,
		PaymentSvc: paymentSvc,
	}
}

func (pc *PaymentControllerImpl) GeneratePaymentInfo(ctx echo.Context) error {
	var paymentInfoReq model.GeneratePaymentInfoRequest

	if err := ctx.Bind(&paymentInfoReq); err != nil {
		return helper.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), err.Error(), err, nil)
	}

	results, err := pc.PaymentSvc.GeneratePaymentInfo(ctx.Request().Context(), &paymentInfoReq)
	if err != nil {
		return helper.NewResponses[any](ctx, http.StatusInternalServerError, "Error Generate Payment Info", nil, err, nil)
	}

	return helper.NewResponses[any](ctx, http.StatusOK, "Success Generate Payment Info", results, nil, nil)
}

func (pc *PaymentControllerImpl) CreatePayment(ctx echo.Context) error {
	var paymentReq model.CreateTransactionRequest

	if err := ctx.Bind(&paymentReq); err != nil {
		return helper.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), err.Error(), err, nil)
	}

	err := paymentReq.Validate()
	if err != nil {
		return helper.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), err.Error(), err, nil)
	}

	results, err := pc.PaymentSvc.CreatePayment(ctx.Request().Context(), &paymentReq)
	if err != nil {
		return helper.NewResponses[any](ctx, http.StatusInternalServerError, "Error Create Payment", nil, err, nil)
	}

	return helper.NewResponses[any](ctx, http.StatusCreated, "Success Create Payment", results, nil, nil)
}

func (pc *PaymentControllerImpl) PaymentNotification(ctx echo.Context) error {
	return nil
}
