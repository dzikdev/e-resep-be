package requester

import (
	"context"
	"e-resep-be/internal/config"
	"e-resep-be/internal/helper"
	"e-resep-be/internal/model"
	"net/http"

	"github.com/sirupsen/logrus"
)

type (
	// KimiaFarmaRequester is an interface that has all the function to be implemented inside kimia farma requester
	KimiaFarmaRequester interface {
		CheckAvailabilityAndPriceMedicationByCode(ctx context.Context, kfaCode string) (*model.CheckAvailabilityResponse, error)
	}

	// KimiaFarmaRequesterImpl is an app kimia farma struct that consists of all the dependencies needed for kimia farma requester
	KimiaFarmaRequesterImpl struct {
		Context    context.Context
		Config     *config.Configuration
		Logger     *logrus.Logger
		HTTPClient *http.Client
	}
)

// NewKimiaFarmaRequester return new instances kimia farma requester
func NewKimiaFarmaRequester(ctx context.Context, config *config.Configuration, logger *logrus.Logger, httpCli *http.Client) *KimiaFarmaRequesterImpl {
	return &KimiaFarmaRequesterImpl{
		Context:    ctx,
		Config:     config,
		Logger:     logger,
		HTTPClient: httpCli,
	}
}

func (kr *KimiaFarmaRequesterImpl) CheckAvailabilityAndPriceMedicationByCode(ctx context.Context, kfaCode string) (*model.CheckAvailabilityResponse, error) {
	// TODO : implement requester check availablilty and price medication by id
	if kr.Config.Server.AppEnv == "development" {
		return &model.CheckAvailabilityResponse{
			IsAvailable: true,
			Price:       helper.GenerateRandomPrice(3000, 200000),
		}, nil
	}
	return nil, nil
}
