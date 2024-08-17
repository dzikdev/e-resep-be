package repository

import (
	"context"
	"e-resep-be/internal/config"
	"e-resep-be/internal/helper"
	"e-resep-be/internal/model"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type (
	// TransactionRepository is an interface that has all the function to be implemented inside transaction repository
	TransactionRepository interface {
		Insert(ctx context.Context, req *model.CreateTransactionRequest) (int, error)
		GetDetailsByTransactionID(ctx context.Context, transactionID int) ([]model.TransactionDetail, error)
		UpdateStatusByID(ctx context.Context, status model.TransactionStatusEnum, id int) error
	}

	// TransactionRepositoryImpl is an app transaction struct that consists of all the dependencies needed for transaction repository
	TransactionRepositoryImpl struct {
		Context context.Context
		Config  *config.Configuration
		Logger  *logrus.Logger
		DB      *pgxpool.Pool
	}
)

// NewTransactionRepository return new instances transaction repository
func NewTransactionRepository(ctx context.Context, config *config.Configuration, logger *logrus.Logger, db *pgxpool.Pool) *TransactionRepositoryImpl {
	return &TransactionRepositoryImpl{
		Context: ctx,
		Config:  config,
		Logger:  logger,
		DB:      db,
	}
}

func (tr *TransactionRepositoryImpl) Insert(ctx context.Context, req *model.CreateTransactionRequest) (int, error) {
	qInsertTrx := `
		INSERT INTO transaction (patient_id, patient_address_id, status, additional_price, total_price) VALUES ($1,$2,$3,$4,$5) RETURNING id
	`

	qInsertTrxDetail := `
		INSERT INTO transaction_detail (transaction_id, medication_id, medication_name, price) VALUES %s
	`

	tx, err := tr.DB.Begin(ctx)
	if err != nil {
		tr.Logger.Error("TransactionRepositoryImpl.Insert Begin TX ERROR", err)

		return 0, err
	}

	var transactionID int
	row := tx.QueryRow(ctx, qInsertTrx, req.PatientID, req.PatientAddressID, model.TransactionStatusEnumPending, req.AdditionalPrice, req.TotalPrice)
	err = row.Scan(
		&transactionID,
	)
	if err != nil {

		errRollback := tx.Rollback(ctx)
		if errRollback != nil {
			tr.Logger.Error("TransactionRepositoryImpl.Insert ERROR rollback TX", errRollback)

			return 0, errRollback
		}
		tr.Logger.Error("TransactionRepositoryImpl.Insert row Scan ERROR", err)

		return 0, err
	}

	numberArgsPerRows := 4
	valueArgs := make([]interface{}, 0, numberArgsPerRows*len(req.Items))

	for i := 0; i < len(req.Items); i++ {
		valueArgs = append(valueArgs, transactionID, req.Items[i].ID, req.Items[i].Name, req.Items[i].Price)
	}

	qInsertTrxDetail = helper.BulkInsert(qInsertTrxDetail, numberArgsPerRows, len(req.Items))

	_, err = tx.Exec(ctx, qInsertTrxDetail, valueArgs...)
	if err != nil {
		errRollback := tx.Rollback(ctx)
		if errRollback != nil {
			tr.Logger.Error("TransactionRepositoryImpl.Insert ERROR rollback TX", errRollback)

			return 0, errRollback
		}

		tr.Logger.Error("TransactionRepositoryImpl.Insert ERROR Exec Insert Bulk Transaction Detail", err)

		return 0, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		errRollback := tx.Rollback(ctx)
		if errRollback != nil {
			tr.Logger.Error("TransactionRepositoryImpl.Insert ERROR rollback TX", errRollback)

			return 0, errRollback
		}

		tr.Logger.Error("TransactionRepositoryImpl.Insert ERROR commit TX", err)

		return 0, err
	}

	return transactionID, nil
}

func (tr *TransactionRepositoryImpl) GetDetailsByTransactionID(ctx context.Context, transactionID int) ([]model.TransactionDetail, error) {
	q := `
		SELECT
			id,
			transaction_id,
			medication_id,
			medication_name,
			price,
			created_at
		FROM
			transaction_detail
		WHERE
			transaction_id = $1
	`

	transactionDetails := []model.TransactionDetail{}

	rows, err := tr.DB.Query(ctx, q, transactionID)
	if err != nil {
		tr.Logger.Error("TransactionRepositoryImpl.GetDetailsByTransactionID Query ERROR", err)

		return []model.TransactionDetail{}, err
	}

	for rows.Next() {
		trxDetail := model.TransactionDetail{}
		err := rows.Scan(
			&trxDetail.ID,
			&trxDetail.TransactionID,
			&trxDetail.MedicationID,
			&trxDetail.MedicationName,
			&trxDetail.Price,
			&trxDetail.CreatedAt,
		)
		if err != nil {
			tr.Logger.Error("TransactionRepositoryImpl.GetDetailsByTransactionID rows Scan ERROR", err)

			return []model.TransactionDetail{}, err
		}

		transactionDetails = append(transactionDetails, trxDetail)
	}

	return transactionDetails, nil
}

func (tr *TransactionRepositoryImpl) UpdateStatusByID(ctx context.Context, status model.TransactionStatusEnum, id int) error {
	q := `
		UPDATE transaction SET status = $1 WHERE id = $2
	`

	_, err := tr.DB.Exec(ctx, q, status, id)
	if err != nil {
		tr.Logger.Error("TransactionRepositoryImpl.UpdateStatusByID Exec ERROR", err)

		return err
	}

	return nil
}
