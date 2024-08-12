package service

import (
	"context"
	"e-resep-be/internal/config"
	"e-resep-be/internal/model"
	"e-resep-be/internal/repository"
)

type (
	// AddressService is an interface that has all the function to be implemented inside address service
	AddressService interface {
		GetProvince(ctx context.Context) (*[]model.Province, error)
		GetCityByProvinceID(ctx context.Context, provinceID int) (*[]model.City, error)
		GetDistrictByCityID(ctx context.Context, cityID int) (*[]model.District, error)
		GetSubDistrictByDistrictID(ctx context.Context, districtID int) (*[]model.SubDistrict, error)
	}

	// AddressServiceImpl is an app address struct that consists of all the dependencies needed for address service
	AddressServiceImpl struct {
		Context     context.Context
		Config      *config.Configuration
		AddressRepo repository.AddressRepository
	}
)

// NewAddressService return new instances address service
func NewAddressService(ctx context.Context, config *config.Configuration, addressRepo repository.AddressRepository) *AddressServiceImpl {
	return &AddressServiceImpl{
		Context:     ctx,
		Config:      config,
		AddressRepo: addressRepo,
	}
}

func (as *AddressServiceImpl) GetProvince(ctx context.Context) (*[]model.Province, error) {
	return as.AddressRepo.GetProvince(ctx)
}

func (as *AddressServiceImpl) GetCityByProvinceID(ctx context.Context, provinceID int) (*[]model.City, error) {
	return as.AddressRepo.GetCityByProvinceID(ctx, provinceID)
}

func (as *AddressServiceImpl) GetDistrictByCityID(ctx context.Context, cityID int) (*[]model.District, error) {
	return as.AddressRepo.GetDistrictByCityID(ctx, cityID)
}

func (as *AddressServiceImpl) GetSubDistrictByDistrictID(ctx context.Context, districtID int) (*[]model.SubDistrict, error) {
	return as.AddressRepo.GetSubDistrictByDistrictID(ctx, districtID)
}
