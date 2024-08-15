package repository

import (
	"context"
	"e-resep-be/internal/config"
	"e-resep-be/internal/model"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type (
	// MedicationRepository is an interface that has all the function to be implemented inside medication repository
	MedicationRepository interface {
		GetByID(ctx context.Context, id int) (*model.MedicationDB, error)
	}

	// MedicationRepositoryImpl is an app medication struct that consists of all the dependencies needed for medication repository
	MedicationRepositoryImpl struct {
		Context context.Context
		Config  *config.Configuration
		Logger  *logrus.Logger
		DB      *pgxpool.Pool
	}
)

// NewMedicationRepository return new instances medication repository
func NewMedicationRepository(ctx context.Context, config *config.Configuration, logger *logrus.Logger, db *pgxpool.Pool) *MedicationRepositoryImpl {
	return &MedicationRepositoryImpl{
		Context: ctx,
		Config:  config,
		Logger:  logger,
		DB:      db,
	}
}

func (mr *MedicationRepositoryImpl) GetByID(ctx context.Context, id int) (*model.MedicationDB, error) {
	q := `
		SELECT
			id,
			ref_id,
			code,
			code_display AS display
		FROM
			medication
		WHERE
			id = $1
	`

	medication := model.MedicationDB{}
	row := mr.DB.QueryRow(ctx, q, id)
	err := row.Scan(
		&medication.ID,
		&medication.RefID,
		&medication.Code,
		&medication.Display,
	)
	if err != nil {
		mr.Logger.Error("MedicationRepositoryImpl.GetByID row Scan ERROR", err)
		return nil, err
	}

	return &medication, nil
}
