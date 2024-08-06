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
		SendMessageByRecipentNumber(ctx context.Context, destination string, templateName string, parameters []model.Parameter, accessToken string) error
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

func (wr *WhatsappRequesterImpl) SendMessageByRecipentNumber(ctx context.Context, destination string, templateName string, parameters []model.Parameter, accessToken string) error {
	msg := model.Message{
		MessagingProduct: "whatsapp",
		To:               destination,
		Type:             "template",
		Template: model.Template{
			Name:     templateName,
			Language: model.Language{Code: "id_ID"},
			Components: []model.Component{
				{
					Type:       "body",
					Parameters: parameters,
				},
			},
		},
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error marshaling message: %v", err)
	}

	req, err := http.NewRequest("POST", wr.Config.Whatsapp.WaBroadcastURL, bytes.NewBuffer(msgBytes))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

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
