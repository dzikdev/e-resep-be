package model

import (
	"time"
)

type Medication struct {
	Code struct {
		Coding []Coding `json:"coding"`
	} `json:"code"`
	Extension []struct {
		URL                  string               `json:"url"`
		ValueCodeableConcept ValueCodeableConcept `json:"valueCodeableConcept"`
	} `json:"extension"`
	Form struct {
		Coding []Coding `json:"coding"`
	} `json:"form"`
	ID           string       `json:"id"`
	Identifier   []Identifier `json:"identifier"`
	Ingredient   []Ingredient `json:"ingredient"`
	Amount       interface{}  `json:"amount"`
	Batch        interface{}  `json:"batch"`
	Manufacturer struct {
		Reference string `json:"reference"`
	} `json:"manufacturer"`
	Meta struct {
		LastUpdated time.Time `json:"lastUpdated"`
		Profile     []string  `json:"profile"`
		VersionID   string    `json:"versionId"`
	} `json:"meta"`
	ResourceType string `json:"resourceType"`
	Status       string `json:"status"`
}

type Coding struct {
	Code    string `json:"code"`
	Display string `json:"display"`
	System  string `json:"system"`
}

type ValueCodeableConcept struct {
	Coding []Coding `json:"coding"`
}

type Identifier struct {
	System string `json:"system"`
	Use    string `json:"use"`
	Value  string `json:"value"`
}

type Ingredient struct {
	IsActive            bool                `json:"isActive"`
	ItemCodeableConcept ItemCodeableConcept `json:"itemCodeableConcept"`
	Strength            Strength            `json:"strength"`
}

type ItemCodeableConcept struct {
	Coding []Coding `json:"coding"`
}

type Strength struct {
	Denominator Quantity `json:"denominator"`
	Numerator   Quantity `json:"numerator"`
}

type Quantity struct {
	Code   string  `json:"code"`
	System string  `json:"system"`
	Value  float64 `json:"value"`
}
