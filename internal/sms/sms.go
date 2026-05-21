package sms

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"notification-service/internal/logger"

	"go.uber.org/zap"
)

func SendOTP(
	phone string,
	otp string,
) error {

	message :=
		otp + " is your OTP for Registration"

	apiURL :=
		fmt.Sprintf(
			"%s?userid=%s&pwd=%s&mobile=971%s&sender=%s&msg=%s&msgtype=16",

			os.Getenv("SMS_API_URL"),

			os.Getenv("SMS_USERID"),

			url.QueryEscape(
				os.Getenv("SMS_PASSWORD"),
			),

			phone,

			os.Getenv("SMS_SENDER"),

			url.QueryEscape(message),
		)

	resp, err := http.Get(apiURL)

	if err != nil {

		logger.Log.Error(
			"sms sending failed",

			zap.Error(err),
		)

		return err
	}

	defer resp.Body.Close()

	logger.Log.Info(
		"sms sent successfully",

		zap.String(
			"phone",
			phone,
		),
	)

	return nil
}
