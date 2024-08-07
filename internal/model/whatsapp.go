package model

type TemplateName string

type SendMessageRequest struct {
	To          string `json:"to"`
	TypeMessage string `json:"type_message"`
	Message     string `json:"message"`
}

const (
	TemplateSendPrescription TemplateName = "SEND_PRESCRIPTION"
)
