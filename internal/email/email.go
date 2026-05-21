package email

import (
	"notification-service/internal/logger"

	"go.uber.org/zap"
)

func SendEmail(
	to string,
	subject string,
	body string,
) error {

	logger.Log.Info(
		"email sent",

		zap.String(
			"to",
			to,
		),

		zap.String(
			"subject",
			subject,
		),
	)

	return nil
}
