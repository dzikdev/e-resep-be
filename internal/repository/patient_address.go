package repository

import (
	"context"
	"e-resep-be/internal/config"
	"e-resep-be/internal/model"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type (
	// PatientAddressRepository is an interface that has all the function to be implemented inside patient address repository
	PatientAddressRepository interface {
		Insert(ctx context.Context, req *model.CreateOrUpdatePatientAddressRequest, patientID int) error
		UpdateByID(ctx context.Context, req *model.CreateOrUpdatePatientAddressRequest, id int) error
	}

	// AddressRepositoryImpl is an app address struct that consists of all the dependencies needed for patient address repository
	PatientAddressRepositoryImpl struct {
		Context context.Context
		Config  *config.Configuration
		Logger  *logrus.Logger
		DB      *pgxpool.Pool
	}
)

// NewPatientAddressRepository return new instances patient address repository
func NewPatientAddressRepository(ctx context.Context, config *config.Configuration, logger *logrus.Logger, db *pgxpool.Pool) *PatientAddressRepositoryImpl {
	return &PatientAddressRepositoryImpl{
		Context: ctx,
		Config:  config,
		Logger:  logger,
		DB:      db,
	}
}

func (pr *PatientAddressRepositoryImpl) Insert(ctx context.Context, req *model.CreateOrUpdatePatientAddressRequest, patientID int) error {
	q := `
		INSERT INTO patient_address (patient_id, address,district, sub_district, city, province, postal_code, coordinates, recipent_name, recipent_phone_number, additional_notes) VALUES ($1,$2,$3,$4,$5,$6,$7,POINT($8,$9),$10,$11,$12)
	`

	_, err := pr.DB.Exec(ctx, q, patientID, req.Address, req.District, req.SubDistrict, req.City, req.Province, req.PostalCode, req.Latitude, req.Longitude, req.RecipentName, req.RecipentPhoneNumber, req.AdditionalNotes)
	if err != nil {
		pr.Logger.Error("PatientAddressRepositoryImpl.Insert ERROR", err)

		return err
	}

	return nil
}

func (pr *PatientAddressRepositoryImpl) UpdateByID(ctx context.Context, req *model.CreateOrUpdatePatientAddressRequest, id int) error {
	q := `
		UPDATE patient_address SET address = $1,district = $2, sub_district = $3, city = $4, province = $5, postal_code = $6, coordinates = POINT($7, $8), recipent_name = $9, recipent_phone_number = $10, additional_notes = $11 WHERE id = $12
	`

	_, err := pr.DB.Exec(ctx, q, req.Address, req.District, req.SubDistrict, req.City, req.Province, req.PostalCode, req.Latitude, req.Longitude, req.RecipentName, req.RecipentPhoneNumber, req.AdditionalNotes, id)
	if err != nil {
		pr.Logger.Error("PatientAddressRepositoryImpl.UpdateByID ERROR", err)

		return err
	}

	return nil
}