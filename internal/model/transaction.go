package model

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	TransactionStatusEnum string

	CreateTransactionRequest struct {
		PatientID        int    `db:"patient_id" json:"patient_id"`
		PatientAddressID int    `db:"patient_address_id" json:"patient_address_id"`
		Items            []Item `json:"items"`
		AdditionalPrice  int    `db:"additional_price" json:"additional_price"`
		TotalPrice       int    `db:"total_price" json:"total_price"`
	}

	TransactionDetail struct {
		ID             int       `db:"id" json:"id"`
		TransactionID  int       `db:"transaction_id" json:"transaction_id"`
		MedicationID   int       `db:"medication_id" json:"medication_id"`
		MedicationName string    `db:"medication_name" json:"medication_name"`
		Price          int       `db:"price" json:"price"`
		CreatedAt      time.Time `db:"created_at" json:"created_at"`
	}
)

const (
	TransactionStatusEnumPending TransactionStatusEnum = "PENDING"
	TransactionStatusEnumProcess TransactionStatusEnum = "PROCESS"
	TransactionStatusEnumSuccess TransactionStatusEnum = "SUCCESS"
	TransactionStatusEnumFailed  TransactionStatusEnum = "FAILED"
	TransactionStatusEnumExpired TransactionStatusEnum = "EXPIRED"
)

func (v CreateTransactionRequest) Validate() error {
	if err := validation.ValidateStruct(&v,
		validation.Field(&v.PatientID, validation.Required),
		validation.Field(&v.PatientAddressID, validation.Required),
		validation.Field(&v.Items, validation.Required),
		validation.Field(&v.TotalPrice, validation.Required),
	); err != nil {
		return err
	}

	return nil
}
