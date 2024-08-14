package service

import (
	"context"
	"e-resep-be/internal/config"
	"e-resep-be/internal/model"
	"e-resep-be/internal/repository"
)

type (
	// PatientAddressService is an interface that has all the function to be implemented inside patient address service
	PatientAddressService interface {
		Create(ctx context.Context, req *model.CreateOrUpdatePatientAddressRequest) error
		Update(ctx context.Context, req *model.CreateOrUpdatePatientAddressRequest, id int) error
	}

	// PatientAddressServiceImpl is an app prescription struct that consists of all the dependencies needed for patient address service
	PatientAddressServiceImpl struct {
		Context            context.Context
		Config             *config.Configuration
		PatientRepo        repository.PatientRepository
		PatientAddressRepo repository.PatientAddressRepository
	}
)

// NewPatientAddressService return new instances patient address service
func NewPatientAddressService(ctx context.Context, config *config.Configuration, patientRepo repository.PatientRepository, patientAddressRepo repository.PatientAddressRepository) *PatientAddressServiceImpl {
	return &PatientAddressServiceImpl{
		Context:            ctx,
		Config:             config,
		PatientRepo:        patientRepo,
		PatientAddressRepo: patientAddressRepo,
	}
}

func (ps *PatientAddressServiceImpl) Create(ctx context.Context, req *model.CreateOrUpdatePatientAddressRequest) error {
	patient, err := ps.PatientRepo.GetByRefID(ctx, req.PatientID)
	if err != nil {
		return err
	}

	err = ps.PatientAddressRepo.Insert(ctx, req, patient.ID)
	if err != nil {
		return err
	}

	return nil
}

func (ps *PatientAddressServiceImpl) Update(ctx context.Context, req *model.CreateOrUpdatePatientAddressRequest, id int) error {
	return ps.PatientAddressRepo.UpdateByID(ctx, req, id)
}
