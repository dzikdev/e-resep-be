package model

type (
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
)
