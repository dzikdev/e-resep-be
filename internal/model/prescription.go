package model

type (
	PrescriptionRequest struct {
		Medication        Medication        `json:"medication"`
		MedicationRequest MedicationRequest `json:"medicationRequest"`
	}

	Prescription struct {
		ID          string `db:"id" json:"id"`
		Display     string `db:"display" json:"display"`
		Code        string `db:"code" json:"code"`
		Price       int    `json:"price"`
		IsAvailable bool   `json:"isAvailable"`
	}
)
