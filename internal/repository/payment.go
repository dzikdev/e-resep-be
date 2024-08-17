package repository

import (
	"context"
	"e-resep-be/internal/config"
	"e-resep-be/internal/model"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type (
	// PaymentRepository is an interface that has all the function to be implemented inside payment repository
	PaymentRepository interface {
		Insert(ctx context.Context, req *model.CreatePaymentRequest) (int, error)
		UpdateStatusByID(ctx context.Context, status model.PaymentStatusEnum, id int) error
		UpdatePartnerIDByID(ctx context.Context, partnerID string, id int) error
	}

	// TransactionRepositoryImpl is an app payment struct that consists of all the dependencies needed for payment repository
	PaymentRepositoryImpl struct {
		Context context.Context
		Config  *config.Configuration
		Logger  *logrus.Logger
		DB      *pgxpool.Pool
	}
)

// NewPaymentRepository return new instances payment repository
func NewPaymentRepository(ctx context.Context, config *config.Configuration, logger *logrus.Logger, db *pgxpool.Pool) *PaymentRepositoryImpl {
	return &PaymentRepositoryImpl{
		Context: ctx,
		Config:  config,
		Logger:  logger,
		DB:      db,
	}
}

func (pr *PaymentRepositoryImpl) Insert(ctx context.Context, req *model.CreatePaymentRequest) (int, error) {
	q := `
		INSERT INTO payment (transaction_id, status, final_price) VALUES ($1,$2,$3) RETURNING id
	`

	var paymentID int
	row := pr.DB.QueryRow(ctx, q, req.TransactionID, model.PaymentStatusEnumPending, req.FinalPrice)
	err := row.Scan(
		&paymentID,
	)
	if err != nil {
		pr.Logger.Error("PaymentRepositoryImpl.Insert QueryRow Scan ERROR", err)

		return 0, err
	}

	return paymentID, nil
}

func (pr *PaymentRepositoryImpl) UpdateStatusByID(ctx context.Context, status model.PaymentStatusEnum, id int) error {
	q := `
		UPDATE payment SET status = $1 WHERE id = $2
	`

	_, err := pr.DB.Exec(ctx, q, status, id)
	if err != nil {
		pr.Logger.Error("PaymentRepositoryImpl.UpdateStatusByID Exec ERROR", err)

		return err
	}

	return nil
}

func (pr *PaymentRepositoryImpl) UpdatePartnerIDByID(ctx context.Context, partnerID string, id int) error {
	q := `
		UPDATE payment SET partner_id = $1 WHERE id = $2
	`

	_, err := pr.DB.Exec(ctx, q, partnerID, id)
	if err != nil {
		pr.Logger.Error("PaymentRepositoryImpl.UpdatePartnerIDByID Exec ERROR", err)

		return err
	}

	return nil
}
