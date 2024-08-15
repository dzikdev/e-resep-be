package service

import (
	"context"
	"e-resep-be/internal/config"
	"e-resep-be/internal/model"
	"e-resep-be/internal/repository"
	"e-resep-be/internal/requester"
	"fmt"
)

type (
	// PaymentService is an interface that has all the function to be implemented inside payment service
	PaymentService interface {
		GeneratePaymentInfo(ctx context.Context, req *model.GeneratePaymentInfoRequest) (*model.PaymentInfo, error)
	}

	// PaymentServiceImpl is an app payment struct that consists of all the dependencies needed for payment service
	PaymentServiceImpl struct {
		Context             context.Context
		Config              *config.Configuration
		MedicationRepo      repository.MedicationRepository
		PatientRepo         repository.PatientRepository
		PatientAddressRepo  repository.PatientAddressRepository
		KimiaFarmaRequester requester.KimiaFarmaRequester
	}
)

// NewPaymentService return new instances payment service
func NewPaymentService(ctx context.Context, config *config.Configuration, medicationRepo repository.MedicationRepository, patientRepo repository.PatientRepository, patientAddressRepo repository.PatientAddressRepository, kimiaFarmaRequester requester.KimiaFarmaRequester) *PaymentServiceImpl {
	return &PaymentServiceImpl{
		Context:             ctx,
		Config:              config,
		MedicationRepo:      medicationRepo,
		PatientRepo:         patientRepo,
		PatientAddressRepo:  patientAddressRepo,
		KimiaFarmaRequester: kimiaFarmaRequester,
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
