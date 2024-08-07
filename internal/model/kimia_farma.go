package model

type (
	// TODO : temporary response mapping
	CheckAvailabilityResponse struct {
		IsAvailable bool `json:"is_available"`
		Price       int  `json:"price"`
	}
)
