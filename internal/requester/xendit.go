package requester

import (
	"context"
	"e-resep-be/internal/config"
	"encoding/json"

	"github.com/sirupsen/logrus"
	"github.com/xendit/xendit-go/v6"
	"github.com/xendit/xendit-go/v6/invoice"
)

type (
	// XenditRequester is an interface that has all the function to be implemented inside xendit requester
	XenditRequester interface {
		CreateInvoice(ctx context.Context, refID string, req invoice.CreateInvoiceRequest) (*invoice.Invoice, error)
		GetInvoiceByID(ctx context.Context, invoiceID string) (*invoice.Invoice, error)
		ExpireInvoiceByID(ctx context.Context, invoiceID string) (*invoice.Invoice, error)
	}

	// XenditRequesterImpl is an app xendit struct that consists of all the dependencies needed for xendit requester
	XenditRequesterImpl struct {
		Context   context.Context
		Config    *config.Configuration
		Logger    *logrus.Logger
		XenditSDK *xendit.APIClient
	}
)

// NewXenditRequester return new instances xendit requester
func NewXenditRequester(ctx context.Context, config *config.Configuration, logger *logrus.Logger, xenditSDK *xendit.APIClient) *XenditRequesterImpl {
	return &XenditRequesterImpl{
		Context:   ctx,
		Config:    config,
		Logger:    logger,
		XenditSDK: xenditSDK,
	}
}

func (xr *XenditRequesterImpl) CreateInvoice(ctx context.Context, refID string, req invoice.CreateInvoiceRequest) (*invoice.Invoice, error) {
	resp, _, err := xr.XenditSDK.InvoiceApi.CreateInvoice(ctx).
		CreateInvoiceRequest(req).
		Execute()
	if err != nil {
		xr.Logger.Error("XenditRequesterImpl.CreateInvoice.Execute ERROR Message ", err.Error())

		b, _ := json.Marshal(err.FullError())
		xr.Logger.Error("XenditRequesterImpl.CreateInvoice.Execute Full Error Struct", string(b))

		return nil, err
	}

	xr.Logger.Info("Success Create Invoice ", resp)

	return resp, nil
}

func (xr *XenditRequesterImpl) GetInvoiceByID(ctx context.Context, invoiceID string) (*invoice.Invoice, error) {
	resp, _, err := xr.XenditSDK.InvoiceApi.GetInvoiceById(ctx, invoiceID).
		Execute()
	if err != nil {
		xr.Logger.Error("XenditRequesterImpl.GetInvoiceByID.Execute ERROR Message ", err.Error())

		b, _ := json.Marshal(err.FullError())
		xr.Logger.Error("XenditRequesterImpl.GetInvoiceByID.Execute Full Error Struct", string(b))

		return nil, err
	}

	xr.Logger.Info("Success Get Invoice ", resp)

	return resp, nil
}

func (xr *XenditRequesterImpl) ExpireInvoiceByID(ctx context.Context, invoiceID string) (*invoice.Invoice, error) {
	resp, _, err := xr.XenditSDK.InvoiceApi.ExpireInvoice(ctx, invoiceID).
		Execute()
	if err != nil {
		xr.Logger.Error("XenditRequesterImpl.ExpireInvoiceByID.Execute ERROR Message ", err.Error())

		b, _ := json.Marshal(err.FullError())
		xr.Logger.Error("XenditRequesterImpl.ExpireInvoiceByID.Execute Full Error Struct", string(b))

		return nil, err
	}

	xr.Logger.Info("Success Expire Manually Invoice ", resp)

	return resp, nil
}
