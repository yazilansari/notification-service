package whatsapp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"notification-service/internal/logger"

	"go.uber.org/zap"
)

func SendWhatsAppOTP(
	phone string,
	otp string,
) error {

	payload := map[string]interface{}{
		"ProfileId": os.Getenv("WHATSAPP_PROFILE_ID"),
		"APIKey":    os.Getenv("WHATSAPP_API_KEY"),

		"MobileNumber": "971" + phone,

		"templateName": "websiteauthentication",

		"Parameters": []string{
			otp,
		},

		"isTemplate": "true",
	}

	jsonData, _ := json.Marshal(payload)

	req, _ := http.NewRequest(
		"POST",
		os.Getenv("WHATSAPP_API_URL"),
		bytes.NewBuffer(jsonData),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {

		logger.Log.Error(
			"whatsapp send failed",

			zap.Error(err),
		)

		return err
	}

	defer resp.Body.Close()

	logger.Log.Info(
		"whatsapp sent successfully",

		zap.String(
			"phone",
			phone,
		),
	)

	return nil
}
