package model

type (
	PrescriptionRequest struct {
		Medication        Medication        `json:"medication"`
		MedicationRequest MedicationRequest `json:"medicationRequest"`
	}
)
