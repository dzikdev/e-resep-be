package model

import (
	"time"
)

type Category struct {
	Coding []Coding `json:"coding"`
}

type CourseOfTherapyType struct {
	Coding []Coding `json:"coding"`
}

type DispenseInterval struct {
	Code   string  `json:"code"`
	System string  `json:"system"`
	Unit   string  `json:"unit"`
	Value  float64 `json:"value"`
}

type ExpectedSupplyDuration struct {
	Code   string  `json:"code"`
	System string  `json:"system"`
	Unit   string  `json:"unit"`
	Value  float64 `json:"value"`
}

type Performer struct {
	Reference string `json:"reference"`
}

type ValidityPeriod struct {
	End   string `json:"end"`
	Start string `json:"start"`
}

type DispenseRequest struct {
	DispenseInterval       DispenseInterval       `json:"dispenseInterval"`
	ExpectedSupplyDuration ExpectedSupplyDuration `json:"expectedSupplyDuration"`
	NumberOfRepeatsAllowed int                    `json:"numberOfRepeatsAllowed"`
	Performer              Performer              `json:"performer"`
	Quantity               Quantity               `json:"quantity"`
	ValidityPeriod         ValidityPeriod         `json:"validityPeriod"`
}

type AdditionalInstruction struct {
	Text string `json:"text"`
}

type DoseQuantity struct {
	Code   string  `json:"code"`
	System string  `json:"system"`
	Unit   string  `json:"unit"`
	Value  float64 `json:"value"`
}

type TypeCoding struct {
	Code    string `json:"code"`
	Display string `json:"display"`
	System  string `json:"system"`
}

type Type struct {
	Coding []TypeCoding `json:"coding"`
}

type DoseAndRate struct {
	DoseQuantity DoseQuantity `json:"doseQuantity"`
	Type         Type         `json:"type"`
}

type Route struct {
	Coding []Coding `json:"coding"`
}

type TimingRepeat struct {
	Frequency  int    `json:"frequency"`
	Period     int    `json:"period"`
	PeriodUnit string `json:"periodUnit"`
}

type Timing struct {
	Repeat TimingRepeat `json:"repeat"`
}

type DosageInstruction struct {
	AdditionalInstruction []AdditionalInstruction `json:"additionalInstruction"`
	DoseAndRate           []DoseAndRate           `json:"doseAndRate"`
	PatientInstruction    string                  `json:"patientInstruction"`
	Route                 Route                   `json:"route"`
	Sequence              int                     `json:"sequence"`
	Text                  string                  `json:"text"`
	Timing                Timing                  `json:"timing"`
}

type Encounter struct {
	Reference string `json:"reference"`
}

type MedicationReference struct {
	Display   string `json:"display"`
	Reference string `json:"reference"`
}

type Meta struct {
	LastUpdated time.Time `json:"lastUpdated"`
	VersionID   string    `json:"versionId"`
}

type ReasonCode struct {
	Coding []Coding `json:"coding"`
}

type Requester struct {
	Display   string `json:"display"`
	Reference string `json:"reference"`
}

type Subject struct {
	Display   string `json:"display"`
	Reference string `json:"reference"`
}

type MedicationRequest struct {
	AuthoredOn          string              `json:"authoredOn"`
	Category            []Category          `json:"category"`
	CourseOfTherapyType CourseOfTherapyType `json:"courseOfTherapyType"`
	DispenseRequest     *DispenseRequest    `json:"dispenseRequest"`
	DosageInstruction   []DosageInstruction `json:"dosageInstruction"`
	Encounter           Encounter           `json:"encounter"`
	ID                  *string             `json:"id"`
	Identifier          []Identifier        `json:"identifier"`
	Intent              string              `json:"intent"`
	MedicationReference MedicationReference `json:"medicationReference"`
	Meta                Meta                `json:"meta"`
	Priority            string              `json:"priority"`
	ReasonCode          []ReasonCode        `json:"reasonCode"`
	Requester           Requester           `json:"requester"`
	ResourceType        string              `json:"resourceType"`
	Status              string              `json:"status"`
	Subject             Subject             `json:"subject"`
}
