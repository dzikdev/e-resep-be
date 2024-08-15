package model

type (
	PrescriptionRequest struct {
		Medication        Medication        `json:"medication"`
		MedicationRequest MedicationRequest `json:"medicationRequest"`
	}

	Prescription struct {
		ID          int    `db:"id" json:"id"`
		Display     string `db:"display" json:"display"`
		Code        string `db:"code" json:"code"`
		PatientID   string `db:"patient_id" json:"patient_id"`
		Price       int    `json:"price"`
		IsAvailable bool   `json:"isAvailable"`
	}
)
