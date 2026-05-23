package whatsapp

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"notification-service/internal/logger"

	"go.uber.org/zap"
)

func normalizePhone(phone string) string {

	phone = strings.TrimSpace(phone)

	phone = strings.TrimPrefix(phone, "+")

	// UAE Format
	// 0501234567 -> 971501234567

	if strings.HasPrefix(phone, "0") {

		phone = "971" + phone[1:]
	}

	return phone
}

func SendWhatsAppOTP(
	phone string,
	otp string,
) error {

	normalizedPhone := normalizePhone(phone)

	payload := map[string]interface{}{
		"ProfileId": os.Getenv("WHATSAPP_PROFILE_ID"),
		"APIKey":    os.Getenv("WHATSAPP_API_KEY"),

		"MobileNumber": normalizedPhone,

		"templateName": "websiteauthentication",

		"Parameters": []string{
			otp,
		},

		"HeaderType": "Text",

		"Text": "",

		"MediaUrl": "",

		"Latitude": 0,

		"Longitude": 0,

		"isTemplate": "true",

		"ButtonOrListJSON": "",

		"SubClientCode": "",

		"HeaderParameter": "",

		"CTAButtonURLParameter": "",

		"CTAButtonURLParameter2": "",
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

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)

	if err != nil {

		logger.Log.Error(
			"whatsapp send failed",

			zap.Error(err),
		)

		return err
	}

	defer resp.Body.Close()

	waBody, _ := io.ReadAll(
		resp.Body,
	)

	logger.Log.Info(
		"whatsapp sent successfully",

		zap.String(
			"phone",
			phone,
		),

		zap.String(
			"WhatsApp API Response:",
			string(waBody),
		),
	)

	return nil
}
