package model

type (
	PatientAdress struct {
		ID        int `db:"id" json:"id"`
		PatientId int `db:"province_id" json:"province_id"`
	}
)
