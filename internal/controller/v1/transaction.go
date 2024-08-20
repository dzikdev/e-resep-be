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
	// TransactionController is an interface that has all the function to be implemented inside transaction controller
	TransactionController interface {
		CreateTransaction(ctx echo.Context) error
		GetTransactionByPartnerID(ctx echo.Context) error
	}

	// TransactionControllerImpl is an app transaction struct that consists of all the dependencies needed for transaction controller
	TransactionControllerImpl struct {
		Context        context.Context
		Config         *config.Configuration
		TransactionSvc service.TransactionService
	}
)

// NewTransactionController return new instance transaction controller
func NewTransactionController(ctx context.Context, config *config.Configuration, transactionSvc service.TransactionService) *TransactionControllerImpl {
	return &TransactionControllerImpl{
		Context:        ctx,
		Config:         config,
		TransactionSvc: transactionSvc,
	}
}

func (tc *TransactionControllerImpl) CreateTransaction(ctx echo.Context) error {
	var transactionReq model.CreateTransactionRequest

	if err := ctx.Bind(&transactionReq); err != nil {
		return helper.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), err.Error(), err, nil)
	}

	err := transactionReq.Validate()
	if err != nil {
		return helper.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), err.Error(), err, nil)
	}

	results, err := tc.TransactionSvc.CreateTransaction(ctx.Request().Context(), &transactionReq)
	if err != nil {
		return helper.NewResponses[any](ctx, http.StatusInternalServerError, "Error Create Transaction", nil, err, nil)
	}

	return helper.NewResponses[any](ctx, http.StatusCreated, "Success Create Transaction", results, nil, nil)
}

func (tc *TransactionControllerImpl) GetTransactionByPartnerID(ctx echo.Context) error {
	results, err := tc.TransactionSvc.CheckStatusByPartnerID(ctx.Request().Context(), ctx.Param("partner_id"))
	if err != nil {
		return helper.NewResponses[any](ctx, http.StatusInternalServerError, "Error Get Transaction By Partner ID", nil, err, nil)
	}

	return helper.NewResponses[any](ctx, http.StatusOK, "Success Get Transaction By Partner ID", results, nil, nil)
}
