package repository

import (
	"context"
	"e-resep-be/internal/config"
	"e-resep-be/internal/model"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type (
	// AddressRepository is an interface that has all the function to be implemented inside address repository
	AddressRepository interface {
		GetProvince(ctx context.Context) (*[]model.Province, error)
		GetCityByProvinceID(ctx context.Context, provinceID int) (*[]model.City, error)
		GetDistrictByCityID(ctx context.Context, cityID int) (*[]model.District, error)
		GetSubDistrictByDistrictID(ctx context.Context, districtID int) (*[]model.SubDistrict, error)
	}

	// AddressRepositoryImpl is an app address struct that consists of all the dependencies needed for address repository
	AddressRepositoryImpl struct {
		Context context.Context
		Config  *config.Configuration
		Logger  *logrus.Logger
		DB      *pgxpool.Pool
	}
)

// NewAddressRepository return new instances address repository
func NewAddressRepository(ctx context.Context, config *config.Configuration, logger *logrus.Logger, db *pgxpool.Pool) *AddressRepositoryImpl {
	return &AddressRepositoryImpl{
		Context: ctx,
		Config:  config,
		Logger:  logger,
		DB:      db,
	}
}

func (ar *AddressRepositoryImpl) GetProvince(ctx context.Context) (*[]model.Province, error) {
	q := `
		SELECT
			id,
			area_id,
			area_name,
			level,
			created_at,
			updated_at,
			delete
		FROM
			province
		WHERE
			delete = 0
		ORDER BY
			area_name ASC
	`

	provinces := []model.Province{}

	rows, err := ar.DB.Query(ctx, q)
	if err != nil {
		ar.Logger.Error("AddressRepositoryImpl.GetProvince.Query ERROR", err)

		return nil, err
	}

	for rows.Next() {
		province := model.Province{}
		err := rows.Scan(
			&province.ID,
			&province.AreaID,
			&province.AreaName,
			&province.Level,
			&province.CreatedAt,
			&province.UpdatedAt,
			&province.Delete,
		)
		if err != nil {
			ar.Logger.Error("AddressRepositoryImpl.GetProvince rows Scan ERROR", err)

			return nil, err
		}

		provinces = append(provinces, province)
	}

	return &provinces, nil
}

func (ar *AddressRepositoryImpl) GetCityByProvinceID(ctx context.Context, provinceID int) (*[]model.City, error) {
	q := `
		SELECT
			id,
			area_id,
			area_name,
			level,
			province_id,
			created_at,
			updated_at,
			delete
		FROM
			city
		WHERE
			delete = 0
		AND
			province_id = $1
		ORDER BY
			area_name ASC
	`

	cities := []model.City{}

	rows, err := ar.DB.Query(ctx, q, provinceID)
	if err != nil {
		ar.Logger.Error("AddressRepositoryImpl.GetCityByProvinceID.Query ERROR", err)

		return nil, err
	}

	for rows.Next() {
		city := model.City{}
		err := rows.Scan(
			&city.ID,
			&city.AreaID,
			&city.AreaName,
			&city.Level,
			&city.ProvinceID,
			&city.CreatedAt,
			&city.UpdatedAt,
			&city.Delete,
		)
		if err != nil {
			ar.Logger.Error("AddressRepositoryImpl.GetCityByProvinceID rows Scan ERROR", err)

			return nil, err
		}

		cities = append(cities, city)
	}

	return &cities, nil
}

func (ar *AddressRepositoryImpl) GetDistrictByCityID(ctx context.Context, cityID int) (*[]model.District, error) {
	q := `
		SELECT
			id,
			area_id,
			area_name,
			level,
			province_id,
			city_id,
			created_at,
			updated_at,
			delete
		FROM
			district
		WHERE
			delete = 0
		AND
			city_id = $1
		ORDER BY
			area_name ASC
	`

	district := []model.District{}

	rows, err := ar.DB.Query(ctx, q, cityID)
	if err != nil {
		ar.Logger.Error("AddressRepositoryImpl.GetDistrictByCityID.Query ERROR", err)

		return nil, err
	}

	for rows.Next() {
		d := model.District{}
		err := rows.Scan(
			&d.ID,
			&d.AreaID,
			&d.AreaName,
			&d.Level,
			&d.ProvinceID,
			&d.CityID,
			&d.CreatedAt,
			&d.UpdatedAt,
			&d.Delete,
		)
		if err != nil {
			ar.Logger.Error("AddressRepositoryImpl.GetDistrictByCityID rows Scan ERROR", err)

			return nil, err
		}

		district = append(district, d)
	}

	return &district, nil
}

func (ar *AddressRepositoryImpl) GetSubDistrictByDistrictID(ctx context.Context, districtID int) (*[]model.SubDistrict, error) {
	q := `
		SELECT
			id,
			area_id,
			area_name,
			level,
			province_id,
			city_id,
			district_id,
			created_at,
			updated_at,
			delete
		FROM
			sub_district
		WHERE
			delete = 0
		AND
			district_id = $1
		ORDER BY
			area_name ASC
	`

	subDistrict := []model.SubDistrict{}

	rows, err := ar.DB.Query(ctx, q, districtID)
	if err != nil {
		ar.Logger.Error("AddressRepositoryImpl.GetSubDistrictByDistrictID.Query ERROR", err)

		return nil, err
	}

	for rows.Next() {
		subD := model.SubDistrict{}
		err := rows.Scan(
			&subD.ID,
			&subD.AreaID,
			&subD.AreaName,
			&subD.Level,
			&subD.ProvinceID,
			&subD.CityID,
			&subD.DistrictID,
			&subD.CreatedAt,
			&subD.UpdatedAt,
			&subD.Delete,
		)
		if err != nil {
			ar.Logger.Error("AddressRepositoryImpl.GetSubDistrictByDistrictID rows Scan ERROR", err)

			return nil, err
		}

		subDistrict = append(subDistrict, subD)
	}

	return &subDistrict, nil
}
