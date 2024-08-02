package service

import (
	"context"

	"e-resep-be/internal/config"
	"e-resep-be/internal/model"
	"e-resep-be/internal/repository"
)

type (
	// PrescriptionService is an interface that has all the function to be implemented inside prescription service
	PrescriptionService interface {
		Create(ctx context.Context, req *model.PrescriptionRequest) error
	}

	// PrescriptionServiceImpl is an app prescription struct that consists of all the dependencies needed for prescription service
	PrescriptionServiceImpl struct {
		Context          context.Context
		Config           *config.Configuration
		PrescriptionRepo repository.PrescriptionRepository
	}
)

// NewPrescriptionService return new instances prescription service
func NewPrescriptionService(ctx context.Context, config *config.Configuration, prescriptionRepo repository.PrescriptionRepository) *PrescriptionServiceImpl {
	return &PrescriptionServiceImpl{
		Context:          ctx,
		Config:           config,
		PrescriptionRepo: prescriptionRepo,
	}
}

func (ps *PrescriptionServiceImpl) Create(ctx context.Context, req *model.PrescriptionRequest) error {
	return ps.PrescriptionRepo.Insert(ctx, req)
}
