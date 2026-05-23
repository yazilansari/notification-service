package sms

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
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

func SendOTP(
	phone string,
	otp string,
) error {

	normalizedPhone := normalizePhone(phone)

	message :=
		otp + " is your OTP for Registration"

	apiURL :=
		fmt.Sprintf(
			"%s?userid=%s&pwd=%s&mobile=%s&sender=%s&msg=%s&msgtype=16",

			os.Getenv("SMS_API_URL"),

			url.QueryEscape(os.Getenv("SMS_USERID")),

			url.QueryEscape(
				os.Getenv("SMS_PASSWORD"),
			),

			url.QueryEscape(normalizedPhone),

			url.QueryEscape(os.Getenv("SMS_SENDER")),

			url.QueryEscape(message),
		)

	smsClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	smsReq, err := http.NewRequest(
		http.MethodGet,
		apiURL,
		nil,
	)

	if err != nil {

		logger.Log.Error(
			"sms sending failed",

			zap.Error(err),
		)

		return err
	}

	smsResp, err := smsClient.Do(
		smsReq,
	)

	if err != nil {

		logger.Log.Error(
			"sms sending failed",

			zap.Error(err),
		)

		return err
	}

	defer smsResp.Body.Close()

	smsBody, _ := io.ReadAll(
		smsResp.Body,
	)

	logger.Log.Info(
		"sms sent successfully",

		zap.String(
			"phone",
			phone,
		),

		zap.String(
			"SMS API Response:",
			string(smsBody),
		),
	)

	return nil
}
