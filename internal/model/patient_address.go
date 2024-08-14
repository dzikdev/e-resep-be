package model

import "time"

type (
	PatientAddress struct {
		ID                  int       `db:"id" json:"id"`
		PatientID           int       `db:"patient_id" json:"patient_id"`
		Address             string    `db:"address" json:"address"`
		SubDistrict         string    `db:"sub_district" json:"sub_district"`
		District            string    `db:"district" json:"district"`
		City                string    `db:"city" json:"city"`
		Province            string    `db:"province" json:"province"`
		PostalCode          string    `db:"postal_code" json:"postal_code"`
		Latitude            float64   `db:"latitude" json:"latitude"`
		Longitude           float64   `db:"longitude" json:"longitude"`
		RecipentName        string    `db:"recipent_name" json:"recipent_name"`
		RecipentPhoneNumber string    `db:"recipent_phone_number" json:"recipent_phone_number"`
		AdditionalNotes     *string   `db:"additional_notes" json:"additional_notes"`
		CreatedAt           time.Time `db:"created_at" json:"created_at"`
	}

	CreateOrUpdatePatientAddressRequest struct {
		Address             string  `db:"address" json:"address"`
		PatientID           string  `json:"patient_id"`
		District            string  `db:"district" json:"district"`
		SubDistrict         string  `db:"sub_district" json:"sub_district"`
		City                string  `db:"city" json:"city"`
		Province            string  `db:"province" json:"province"`
		PostalCode          string  `db:"postal_code" json:"postal_code"`
		Latitude            float64 `db:"latitude" json:"latitude"`
		Longitude           float64 `db:"longitude" json:"longitude"`
		RecipentName        string  `db:"recipent_name" json:"recipent_name"`
		RecipentPhoneNumber string  `db:"recipent_phone_number" json:"recipent_phone_number"`
		AdditionalNotes     string  `db:"additional_notes" json:"additional_notes"`
	}
)
