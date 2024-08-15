package service

import (
	"context"

	"e-resep-be/internal/config"
	"e-resep-be/internal/model"
	"e-resep-be/internal/repository"
	"e-resep-be/internal/requester"
)

type (
	// PrescriptionService is an interface that has all the function to be implemented inside prescription service
	PrescriptionService interface {
		Create(ctx context.Context, req *model.PrescriptionRequest, phoneNumber string) error
		GetByPrescriptionID(ctx context.Context, id string) ([]model.Prescription, error)
	}

	// PrescriptionServiceImpl is an app prescription struct that consists of all the dependencies needed for prescription service
	PrescriptionServiceImpl struct {
		Context             context.Context
		Config              *config.Configuration
		PrescriptionRepo    repository.PrescriptionRepository
		WhatsappRequester   requester.WhatsappRequester
		KimiaFarmaRequester requester.KimiaFarmaRequester
	}
)

// NewPrescriptionService return new instances prescription service
func NewPrescriptionService(ctx context.Context, config *config.Configuration, prescriptionRepo repository.PrescriptionRepository, whatsappRequester requester.WhatsappRequester, kimiaFarmaRequester requester.KimiaFarmaRequester) *PrescriptionServiceImpl {
	return &PrescriptionServiceImpl{
		Context:             ctx,
		Config:              config,
		PrescriptionRepo:    prescriptionRepo,
		WhatsappRequester:   whatsappRequester,
		KimiaFarmaRequester: kimiaFarmaRequester,
	}
}

func (ps *PrescriptionServiceImpl) Create(ctx context.Context, req *model.PrescriptionRequest, phoneNumber string) error {
	// insert to database
	err := ps.PrescriptionRepo.Insert(ctx, req, phoneNumber)
	if err != nil {
		return err
	}

	if phoneNumber != "" {
		// send message to patient number through whatsapp
		err = ps.WhatsappRequester.SendMessageByRecipentNumber(ctx, req.MedicationRequest.Subject.Display, req.MedicationRequest.Identifier[0].Value, phoneNumber, model.TemplateSendPrescription)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ps *PrescriptionServiceImpl) GetByPrescriptionID(ctx context.Context, id string) ([]model.Prescription, error) {
	prescriptions, err := ps.PrescriptionRepo.GetByPrescriptionID(ctx, id)
	if err != nil {
		return []model.Prescription{}, err
	}

	if len(prescriptions) > 0 {
		for i := range prescriptions {
			kimiaFarmaResp, err := ps.KimiaFarmaRequester.CheckAvailabilityAndPriceMedicationByCode(ctx, prescriptions[i].Code)
			if err != nil {
				return []model.Prescription{}, err
			}

			prescriptions[i].IsAvailable = kimiaFarmaResp.IsAvailable
			prescriptions[i].Price = kimiaFarmaResp.Price
		}
	}

	return prescriptions, nil
}
