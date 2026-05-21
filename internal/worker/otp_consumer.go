package worker

import (
	"context"
	"encoding/json"
	"os"

	"notification-service/internal/kafka"
	"notification-service/internal/logger"
	"notification-service/internal/sms"

	kafkaGo "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func StartOTPConsumer() {

	reader := kafkaGo.NewReader(
		kafkaGo.ReaderConfig{
			Brokers: []string{
				os.Getenv("KAFKA_BROKER"),
			},

			Topic: kafka.OTPTopic,

			GroupID: "otp-worker-group",
		},
	)

	logger.Log.Info(
		"otp consumer started",
	)

	for {

		message, err :=
			reader.ReadMessage(
				context.Background(),
			)

		if err != nil {

			logger.Log.Error(
				"failed reading otp message",

				zap.Error(err),
			)

			continue
		}

		var payload kafka.OTPPayload

		err = json.Unmarshal(
			message.Value,
			&payload,
		)

		if err != nil {

			logger.Log.Error(
				"failed parsing otp payload",

				zap.Error(err),
			)

			continue
		}

		err = sms.SendOTP(
			payload.Phone,
			payload.OTP,
		)

		if err != nil {

			logger.Log.Error(
				"failed sending sms",

				zap.Error(err),
			)
		}
	}
}
