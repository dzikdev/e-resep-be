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
	// PaymentService is an interface that has all the function to be implemented inside payment service
	PaymentService interface {
		GeneratePaymentInfo(ctx context.Context, req *model.GeneratePaymentInfoRequest) (*model.PaymentInfo, error)
		CreatePayment(ctx context.Context, req *model.CreateTransactionRequest) (*model.CreatePaymentResponse, error)
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
	totalPriceInItems := 0
	for _, item := range req.Items {
		totalPriceInItems += item.Price
	}

	if totalPriceInItems != req.TotalPrice {
		return nil, model.NewError(model.Validation, "invalid total price")
	}

	transactionID, err := ps.TransactionRepo.Insert(ctx, req)
	if err != nil {
		return nil, err
	}

	getTransactionDetails, err := ps.TransactionRepo.GetDetailsByTransactionID(ctx, transactionID)
	if err != nil {
		return nil, err
	}

	paymentID, err := ps.PaymentRepo.Insert(ctx, &model.CreatePaymentRequest{
		TransactionID: transactionID,
		FinalPrice:    req.TotalPrice,
	})
	if err != nil {
		return nil, err
	}

	invoiceReq := invoice.NewCreateInvoiceRequest(fmt.Sprintf("%d", paymentID), float64(req.TotalPrice))
	invoiceReq.SetDescription(fmt.Sprintf("Create Invoice Transaction E-RESEP with id %d", paymentID))

	for _, trxDetail := range getTransactionDetails {
		item := invoice.InvoiceItem{}

		item.SetReferenceId(fmt.Sprintf("%d", trxDetail.ID))
		item.SetName(trxDetail.MedicationName)
		item.SetPrice(float32(trxDetail.Price))
		item.SetQuantity(1)

		invoiceReq.Items = append(invoiceReq.Items, item)
	}

	results, err := ps.XenditRequester.CreateInvoice(ctx, *invoiceReq)
	if err != nil {
		return nil, err
	}

	err = ps.TransactionRepo.UpdateStatusByID(ctx, model.TransactionStatusEnumProcess, transactionID)
	if err != nil {
		return nil, err
	}

	err = ps.PaymentRepo.UpdateStatusByID(ctx, model.PaymentStatusEnumProcess, paymentID)
	if err != nil {
		return nil, err
	}

	err = ps.PaymentRepo.UpdatePartnerIDByID(ctx, *results.Id, paymentID)
	if err != nil {
		return nil, err
	}

	return &model.CreatePaymentResponse{
		ID:         *results.Id,
		InvoiceURL: results.InvoiceUrl,
	}, nil
}
