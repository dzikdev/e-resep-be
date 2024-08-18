package repository

import (
	"context"
	"e-resep-be/internal/config"
	"e-resep-be/internal/model"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type (
	// PatientRepository is an interface that has all the function to be implemented inside patient repository
	PatientRepository interface {
		GetByRefID(ctx context.Context, refID string) (*model.Patient, error)
		GetByID(ctx context.Context, id int) (*model.Patient, error)
	}

	// PatientRepositoryImpl is an app patient struct that consists of all the dependencies needed for patient repository
	PatientRepositoryImpl struct {
		Context context.Context
		Config  *config.Configuration
		Logger  *logrus.Logger
		DB      *pgxpool.Pool
	}
)

// NewPatientRepository return new instances patient repository
func NewPatientRepository(ctx context.Context, config *config.Configuration, logger *logrus.Logger, db *pgxpool.Pool) *PatientRepositoryImpl {
	return &PatientRepositoryImpl{
		Context: ctx,
		Config:  config,
		Logger:  logger,
		DB:      db,
	}
}

func (pr *PatientRepositoryImpl) GetByRefID(ctx context.Context, refID string) (*model.Patient, error) {
	q := `
		SELECT
			id,
			ref_id,
			name,
			phone_number,
			created_at
		FROM
			patient
		WHERE
			ref_id = $1
	`

	patient := model.Patient{}
	row := pr.DB.QueryRow(ctx, q, refID)
	err := row.Scan(
		&patient.ID,
		&patient.RefID,
		&patient.Name,
		&patient.PhoneNumber,
		&patient.CreatedAt,
	)
	if err != nil {
		pr.Logger.Error("PatientRepositoryImpl.GetByRefID QueryRow.Scan ERROR", err)

		return nil, err
	}

	return &patient, nil
}

func (pr *PatientRepositoryImpl) GetByID(ctx context.Context, id int) (*model.Patient, error) {
	q := `
		SELECT
			id,
			ref_id,
			name,
			phone_number,
			created_at
		FROM
			patient
		WHERE
			id = $1
	`

	patient := model.Patient{}
	row := pr.DB.QueryRow(ctx, q, id)
	err := row.Scan(
		&patient.ID,
		&patient.RefID,
		&patient.Name,
		&patient.PhoneNumber,
		&patient.CreatedAt,
	)
	if err != nil {
		pr.Logger.Error("PatientRepositoryImpl.GetByID QueryRow.Scan ERROR", err)

		return nil, err
	}

	return &patient, nil
}
