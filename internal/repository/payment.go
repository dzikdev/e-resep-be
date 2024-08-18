package repository

import (
	"context"
	"e-resep-be/internal/config"
	"e-resep-be/internal/helper"
	"e-resep-be/internal/model"
	"fmt"
	"reflect"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type (
	// PaymentRepository is an interface that has all the function to be implemented inside payment repository
	PaymentRepository interface {
		Insert(ctx context.Context, req *model.CreatePaymentRequest) (int, error)
		UpdateByID(ctx context.Context, req model.Payment, id int) error
		GetByID(ctx context.Context, id int) (*model.Payment, error)
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

func (pr *PaymentRepositoryImpl) UpdateByID(ctx context.Context, req model.Payment, id int) error {
	q := `
		UPDATE payment SET updated_at = NOW()
	`
	values := make([]interface{}, 0)

	paymentType := reflect.TypeOf(req)
	paymentValue := reflect.ValueOf(req)

	for i := 0; i < paymentType.NumField(); i++ {
		field := paymentType.Field(i)
		value := paymentValue.Field(i)

		if value.Interface() != reflect.Zero(field.Type).Interface() {
			fieldName := field.Tag.Get("db")
			if fieldName == "completed_at" {
				formattedPaidAt := req.CompletedAt.In(helper.TimezoneJakarta).Format("2006-01-02 15:04:05.999999-07:00")
				q += fmt.Sprintf(`, "%s"='%v'`, fieldName, formattedPaidAt)
			} else {
				q += fmt.Sprintf(`, "%s"='%v'`, fieldName, value.Interface())
			}

		}
	}

	q += fmt.Sprintf(` WHERE id='%d'`, id)

	_, err := pr.DB.Exec(ctx, q, values...)
	if err != nil {
		pr.Logger.Error("PaymentRepositoryImpl.UpdateByID Exec ERROR", err)

		return err
	}

	return nil
}

func (pr *PaymentRepositoryImpl) GetByID(ctx context.Context, id int) (*model.Payment, error) {
	q := `
		SELECT
			id,
			transaction_id,
			partner_id,
			completed_at,
			status,
			final_price,
			created_at,
			updated_at
		FROM
			payment
		WHERE
			id = $1
	`

	payment := model.Payment{}
	row := pr.DB.QueryRow(ctx, q, id)
	err := row.Scan(
		&payment.ID,
		&payment.TransactionID,
		&payment.PartnerID,
		&payment.CompletedAt,
		&payment.Status,
		&payment.FinalPrice,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		pr.Logger.Error("PaymentRepositoryImpl.GetByID QueryRow.Scan ERROR", err)

		return nil, err
	}

	return &payment, nil
}
