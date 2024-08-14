package model

import "time"

type (
	Patient struct {
		ID          int       `db:"id" json:"id"`
		RefID       string    `db:"ref_id" json:"ref_id"`
		Name        string    `db:"name" json:"name"`
		PhoneNumber string    `db:"phone_number" json:"phone_number"`
		CreatedAt   time.Time `db:"created_at" json:"created_at"`
	}
)
