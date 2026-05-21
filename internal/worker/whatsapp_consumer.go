package worker

import (
	"context"
	"encoding/json"
	"os"

	"notification-service/internal/kafka"
	"notification-service/internal/logger"
	"notification-service/internal/whatsapp"

	kafkaGo "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func StartWhatsAppConsumer() {

	reader := kafkaGo.NewReader(
		kafkaGo.ReaderConfig{
			Brokers: []string{
				os.Getenv("KAFKA_BROKER"),
			},

			Topic: kafka.WhatsAppTopic,

			GroupID: "whatsapp-worker-group",
		},
	)

	logger.Log.Info(
		"whatsapp consumer started",
	)

	for {

		message, err :=
			reader.ReadMessage(
				context.Background(),
			)

		if err != nil {

			logger.Log.Error(
				"failed reading whatsapp message",

				zap.Error(err),
			)

			continue
		}

		var payload kafka.WhatsAppPayload

		err = json.Unmarshal(
			message.Value,
			&payload,
		)

		if err != nil {

			logger.Log.Error(
				"failed parsing whatsapp payload",

				zap.Error(err),
			)

			continue
		}

		err = whatsapp.SendWhatsAppOTP(
			payload.Phone,
			payload.OTP,
		)

		if err != nil {

			logger.Log.Error(
				"failed sending whatsapp",

				zap.Error(err),
			)
		}
	}
}
