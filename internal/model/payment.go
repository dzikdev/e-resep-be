package model

type (
	PaymentStatusEnum string

	GeneratePaymentInfoRequest struct {
		SelectedMedications []SelectedMedication `json:"selected_medications"`
		PatientID           string               `json:"patient_id"`
	}

	SelectedMedication struct {
		MedicationID int `json:"medication_id"`
	}

	Item struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Price int    `json:"price"`
	}

	PaymentInfo struct {
		PatientAddress *[]PatientAddress `json:"patient_address"`
		Items          []Item            `json:"items"`
		ShippingCost   int               `json:"shipping_cost"`
		TotalPrice     int               `json:"total_price"`
	}

	CreatePaymentRequest struct {
		TransactionID int `db:"transaction_id" json:"transaction_id"`
		FinalPrice    int `db:"final_price" json:"final_price"`
	}

	CreatePaymentResponse struct {
		ID         string `json:"id"`
		InvoiceURL string `json:"invoice_url"`
	}
)

const (
	PaymentStatusEnumPending PaymentStatusEnum = "PENDING"
	PaymentStatusEnumProcess PaymentStatusEnum = "PROCESS"
	PaymentStatusEnumSuccess PaymentStatusEnum = "SUCCESS"
	PaymentStatusEnumFailed  PaymentStatusEnum = "FAILED"
	PaymentStatusEnumExpired PaymentStatusEnum = "EXPIRED"
)
