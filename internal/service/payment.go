package service

import (
	"context"
	"e-resep-be/internal/config"
	"e-resep-be/internal/model"
	"e-resep-be/internal/repository"
	"e-resep-be/internal/requester"
	"fmt"
	"strconv"
	"time"

	"github.com/xendit/xendit-go/v6/invoice"
)

type (
	// PaymentService is an interface that has all the function to be implemented inside payment service
	PaymentService interface {
		GeneratePaymentInfo(ctx context.Context, req *model.GeneratePaymentInfoRequest) (*model.PaymentInfo, error)
		CreatePayment(ctx context.Context, req *model.CreateTransactionRequest) (*model.CreatePaymentResponse, error)
		HandleWebhookNotification(ctx context.Context, req invoice.InvoiceCallback) error
	}

	// PaymentServiceImpl is an app payment struct that consists of all the dependencies needed for payment service
	PaymentServiceImpl struct {
		Context             context.Context
		Config              *config.Configuration
		MedicationRepo      repository.MedicationRepository
		PatientRepo         repository.PatientRepository
		PatientAddressRepo  repository.PatientAddressRepository
		TransactionRepo     repository.TransactionRepository
		PaymentRepo         repository.PaymentRepository
		KimiaFarmaRequester requester.KimiaFarmaRequester
		XenditRequester     requester.XenditRequester
	}
)

// NewPaymentService return new instances payment service
func NewPaymentService(ctx context.Context, config *config.Configuration, medicationRepo repository.MedicationRepository, patientRepo repository.PatientRepository, patientAddressRepo repository.PatientAddressRepository, transactionRepo repository.TransactionRepository, paymentRepo repository.PaymentRepository, kimiaFarmaRequester requester.KimiaFarmaRequester, xenditRequester requester.XenditRequester) *PaymentServiceImpl {
	return &PaymentServiceImpl{
		Context:             ctx,
		Config:              config,
		MedicationRepo:      medicationRepo,
		PatientRepo:         patientRepo,
		PatientAddressRepo:  patientAddressRepo,
		TransactionRepo:     transactionRepo,
		PaymentRepo:         paymentRepo,
		KimiaFarmaRequester: kimiaFarmaRequester,
		XenditRequester:     xenditRequester,
	}
}

func (ps *PaymentServiceImpl) GeneratePaymentInfo(ctx context.Context, req *model.GeneratePaymentInfoRequest) (*model.PaymentInfo, error) {
	resp := model.PaymentInfo{
		Items:          []model.Item{},
		PatientAddress: &[]model.PatientAddress{},
	}

	patient, err := ps.PatientRepo.GetByRefID(ctx, req.PatientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get patient by ref ID: %w", err)
	}

	addresses, err := ps.PatientAddressRepo.GetByPatientID(ctx, patient.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get patient address by patient ID: %w", err)
	}

	if len(*addresses) > 0 {
		for _, addr := range *addresses {
			*resp.PatientAddress = append(*resp.PatientAddress, model.PatientAddress{
				ID:                  addr.ID,
				PatientID:           addr.PatientID,
				Address:             addr.Address,
				District:            addr.District,
				SubDistrict:         addr.SubDistrict,
				City:                addr.City,
				Province:            addr.Province,
				PostalCode:          addr.PostalCode,
				Latitude:            addr.Latitude,
				Longitude:           addr.Longitude,
				RecipentName:        addr.RecipentName,
				RecipentPhoneNumber: addr.RecipentPhoneNumber,
				AdditionalNotes:     addr.AdditionalNotes,
				CreatedAt:           addr.CreatedAt,
			})
		}
	}

	for _, m := range req.SelectedMedications {
		item := model.Item{}

		medication, err := ps.MedicationRepo.GetByID(ctx, m.MedicationID)
		if err != nil {
			return nil, fmt.Errorf("failed to get medication by ID: %w", err)
		}

		medicationDetail, err := ps.KimiaFarmaRequester.CheckAvailabilityAndPriceMedicationByCode(ctx, medication.Code)
		if err != nil {
			return nil, fmt.Errorf("failed to check medication availability and price: %w", err)
		}

		item.ID = medication.ID
		item.Name = medication.Display
		item.Price = medicationDetail.Price

		resp.TotalPrice += medicationDetail.Price
		resp.Items = append(resp.Items, item)
	}

	return &resp, nil
}

func (ps *PaymentServiceImpl) CreatePayment(ctx context.Context, req *model.CreateTransactionRequest) (*model.CreatePaymentResponse, error) {
	// validate total price is same with sum price inside each item
	totalPriceInItems := 0
	for _, item := range req.Items {
		totalPriceInItems += item.Price
	}

	if totalPriceInItems != req.TotalPrice {
		return nil, model.NewError(model.Validation, "invalid total price")
	}

	// get patient by id
	patient, err := ps.PatientRepo.GetByID(ctx, req.PatientID)
	if err != nil {
		return nil, err
	}

	// insert trx & trx details
	transactionID, err := ps.TransactionRepo.Insert(ctx, req)
	if err != nil {
		return nil, err
	}

	// get details transaction by trx id
	getTransactionDetails, err := ps.TransactionRepo.GetDetailsByTransactionID(ctx, transactionID)
	if err != nil {
		return nil, err
	}

	// insert payment
	paymentID, err := ps.PaymentRepo.Insert(ctx, &model.CreatePaymentRequest{
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
	results, err := ps.XenditRequester.CreateInvoice(ctx, *invoiceReq)
	if err != nil {
		return nil, err
	}

	// update transaction status to process by id
	err = ps.TransactionRepo.UpdateByID(ctx, model.Transaction{
		Status: model.TransactionStatusEnumProcess,
	}, transactionID)
	if err != nil {
		return nil, err
	}

	// update payment status to process and fill partner id by id
	err = ps.PaymentRepo.UpdateByID(ctx, model.Payment{
		Status:    model.PaymentStatusEnumProcess,
		PartnerID: *results.Id,
	}, paymentID)
	if err != nil {
		return nil, err
	}

	return &model.CreatePaymentResponse{
		ID:         *results.Id,
		InvoiceURL: results.InvoiceUrl,
	}, nil
}

func (ps *PaymentServiceImpl) HandleWebhookNotification(ctx context.Context, req invoice.InvoiceCallback) error {
	// parse externalId to integer, due payment id is a integer
	parsePaymentID, err := strconv.Atoi(req.ExternalId)
	if err != nil {
		return err
	}

	// recheck payment by id
	payment, err := ps.PaymentRepo.GetByID(ctx, parsePaymentID)
	if err != nil {
		return err
	}

	// recheck transaction by id
	transaction, err := ps.TransactionRepo.GetByID(ctx, payment.TransactionID)
	if err != nil {
		return err
	}

	// update status transaction and payment based on status from xendit callback notification
	switch req.Status {
	case string(invoice.INVOICESTATUS_PENDING):
	case string(invoice.INVOICESTATUS_PAID):
		err := ps.TransactionRepo.UpdateByID(ctx, model.Transaction{
			Status: model.TransactionStatusEnumSuccess,
		}, transaction.ID)
		if err != nil {
			return err
		}

		parsePaidAt, err := time.Parse(time.RFC3339, *req.PaidAt)
		if err != nil {
			return err
		}

		err = ps.PaymentRepo.UpdateByID(ctx, model.Payment{
			Status:      model.PaymentStatusEnumSuccess,
			CompletedAt: &parsePaidAt,
		}, parsePaymentID)
		if err != nil {
			return err
		}
	case string(invoice.INVOICESTATUS_EXPIRED):
		err := ps.TransactionRepo.UpdateByID(ctx, model.Transaction{
			Status: model.TransactionStatusEnumExpired,
		}, transaction.ID)
		if err != nil {
			return err
		}

		err = ps.PaymentRepo.UpdateByID(ctx, model.Payment{
			Status: model.PaymentStatusEnumExpired,
		}, parsePaymentID)
		if err != nil {
			return err
		}
	case string(invoice.INVOICESTATUS_XENDIT_ENUM_DEFAULT_FALLBACK):
		err := ps.TransactionRepo.UpdateByID(ctx, model.Transaction{
			Status: model.TransactionStatusEnumFailed,
		}, transaction.ID)
		if err != nil {
			return err
		}

		err = ps.PaymentRepo.UpdateByID(ctx, model.Payment{
			Status: model.PaymentStatusEnumFailed,
		}, parsePaymentID)
		if err != nil {
			return err
		}
	}

	return nil
}
