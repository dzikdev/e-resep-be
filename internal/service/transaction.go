package service

import (
	"context"
	"e-resep-be/internal/config"
	"e-resep-be/internal/model"
	"e-resep-be/internal/repository"
	"e-resep-be/internal/requester"
	"fmt"

	"github.com/xendit/xendit-go/v6/invoice"
)

type (
	// TransactionService is an interface that has all the function to be implemented inside transaction service
	TransactionService interface {
		CreateTransaction(ctx context.Context, req *model.CreateTransactionRequest) (*model.CreateTransactionResponse, error)
		CheckStatusByPartnerID(ctx context.Context, partnerID string) (*model.CheckStatusTransactionResponse, error)
	}

	// TransactionServiceImpl is an app transaction struct that consists of all the dependencies needed for transaction service
	TransactionServiceImpl struct {
		Context         context.Context
		Config          *config.Configuration
		PatientRepo     repository.PatientRepository
		TransactionRepo repository.TransactionRepository
		PaymentRepo     repository.PaymentRepository
		XenditRequester requester.XenditRequester
	}
)

// NewTransactionService return new instances transaction service
func NewTransactionService(ctx context.Context, config *config.Configuration, patientRepo repository.PatientRepository, transactionRepo repository.TransactionRepository, paymentRepo repository.PaymentRepository, xenditRequester requester.XenditRequester) *TransactionServiceImpl {
	return &TransactionServiceImpl{
		Context:         ctx,
		Config:          config,
		PatientRepo:     patientRepo,
		TransactionRepo: transactionRepo,
		PaymentRepo:     paymentRepo,
		XenditRequester: xenditRequester,
	}
}

func (ts *TransactionServiceImpl) CreateTransaction(ctx context.Context, req *model.CreateTransactionRequest) (*model.CreateTransactionResponse, error) {
	// validate total price is same with sum price inside each item
	totalPriceInItems := 0
	for _, item := range req.Items {
		totalPriceInItems += item.Price
	}

	if totalPriceInItems != req.TotalPrice {
		return nil, model.NewError(model.Validation, "invalid total price")
	}

	// get patient by id
	patient, err := ts.PatientRepo.GetByID(ctx, req.PatientID)
	if err != nil {
		return nil, err
	}

	// insert trx & trx details
	transactionID, err := ts.TransactionRepo.Insert(ctx, req)
	if err != nil {
		return nil, err
	}

	// get details transaction by trx id
	getTransactionDetails, err := ts.TransactionRepo.GetDetailsByTransactionID(ctx, transactionID)
	if err != nil {
		return nil, err
	}

	// insert payment
	paymentID, err := ts.PaymentRepo.Insert(ctx, &model.CreatePaymentRequest{
		TransactionID: transactionID,
		FinalPrice:    req.TotalPrice,
	})
	if err != nil {
		return nil, err
	}

	// initiate invoice object request
	invoiceReq := invoice.NewCreateInvoiceRequest(fmt.Sprintf("%d", paymentID), float64(req.TotalPrice))

	// set customer data inside invoices
	customerData := invoice.NewCustomerObject()
	customerData.SetCustomerId(fmt.Sprintf("%d", patient.ID))
	customerData.SetGivenNames(patient.Name)
	invoiceReq.SetCustomer(*customerData)

	// set description inside invoices
	invoiceReq.SetDescription(fmt.Sprintf("Create Invoice Transaction E-RESEP with id %d", paymentID))

	// set transaction items inside invoices
	for _, trxDetail := range getTransactionDetails {
		item := invoice.NewInvoiceItem(trxDetail.MedicationName, float32(trxDetail.Price), 1)
		item.SetReferenceId(fmt.Sprintf("%d", trxDetail.ID))

		invoiceReq.Items = append(invoiceReq.Items, *item)
	}

	// call xendit to create invoices
	results, err := ts.XenditRequester.CreateInvoice(ctx, *invoiceReq)
	if err != nil {
		return nil, err
	}

	// update transaction status to process by id
	err = ts.TransactionRepo.UpdateByID(ctx, model.Transaction{
		Status: model.TransactionStatusEnumProcess,
	}, transactionID)
	if err != nil {
		return nil, err
	}

	// update payment status to process and fill partner id by id
	err = ts.PaymentRepo.UpdateByID(ctx, model.Payment{
		Status:    model.PaymentStatusEnumProcess,
		PartnerID: *results.Id,
	}, paymentID)
	if err != nil {
		return nil, err
	}

	return &model.CreateTransactionResponse{
		ID:         *results.Id,
		InvoiceURL: results.InvoiceUrl,
	}, nil
}

func (ts *TransactionServiceImpl) CheckStatusByPartnerID(ctx context.Context, partnerID string) (*model.CheckStatusTransactionResponse, error) {
	resp := model.CheckStatusTransactionResponse{
		Items: &[]model.TransactionDetail{},
	}

	payment, err := ts.PaymentRepo.GetByPartnerID(ctx, partnerID)
	if err != nil {
		return nil, err
	}

	transaction, err := ts.TransactionRepo.GetByID(ctx, payment.TransactionID)
	if err != nil {
		return nil, err
	}

	details, err := ts.TransactionRepo.GetDetailsByTransactionID(ctx, payment.TransactionID)
	if err != nil {
		return nil, err
	}

	resp.Transaction = transaction
	*resp.Items = append(*resp.Items, details...)

	return &resp, nil
}
