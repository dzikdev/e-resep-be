package requester

import (
	"bytes"
	"context"
	"e-resep-be/internal/config"
	"e-resep-be/internal/model"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

type (
	// WhatsappRequester is an interface that has all the function to be implemented inside whatsapp requester
	WhatsappRequester interface {
		SendMessageByRecipentNumber(ctx context.Context, patientName, patientID, destination string, templateName model.TemplateName) error
	}

	// WhatsappRequesterImpl is an app whatsapp struct that consists of all the dependencies needed for whatsapp requester
	WhatsappRequesterImpl struct {
		Context    context.Context
		Config     *config.Configuration
		Logger     *logrus.Logger
		HTTPClient *http.Client
	}
)

// NewWhatsappRequester return new instances whatsapp requester
func NewWhatsappRequester(ctx context.Context, config *config.Configuration, logger *logrus.Logger, httpCli *http.Client) *WhatsappRequesterImpl {
	return &WhatsappRequesterImpl{
		Context:    ctx,
		Config:     config,
		Logger:     logger,
		HTTPClient: httpCli,
	}
}

func (wr *WhatsappRequesterImpl) SendMessageByRecipentNumber(ctx context.Context, patientName, patientID, destination string, templateName model.TemplateName) error {

	msgWithTemplate := wr.generateTemplateWhatsappByName(templateName)
	linkFEURL := fmt.Sprintf("%s/resep/%s", wr.Config.Const.ClientURL, patientID)

	sendMessageReq := model.SendMessageRequest{
		To:          destination,
		TypeMessage: "text",
		Message:     fmt.Sprintf(msgWithTemplate, patientName, linkFEURL),
	}

	sendMesssageReqBytes, err := json.Marshal(sendMessageReq)
	if err != nil {
		return fmt.Errorf("error marshaling message: %v", err)
	}

	req, err := http.NewRequest("POST", wr.Config.Whatsapp.WaBroadcastURL, bytes.NewBuffer(sendMesssageReqBytes))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response: %v", resp.Status)
	}

	return nil
}

func (wr *WhatsappRequesterImpl) generateTemplateWhatsappByName(name model.TemplateName) string {
	switch name {
	case model.TemplateSendPrescription:
		return `Halo *%s*,

Resep Anda sudah siap untuk diperiksa. Silakan ikuti tautan di bawah ini untuk mengkonfirmasi ketersediaan dan detail resep Anda:

Periksa Resep Anda di : %s

Jika Anda memiliki pertanyaan, jangan ragu untuk menghubungi tim dukungan kami.

Terima kasih,
*E-RESEP*`
	default:
		return ""
	}
}
