package worker

import (
	"context"
	"encoding/json"
	"os"

	"notification-service/internal/email"
	"notification-service/internal/kafka"
	"notification-service/internal/logger"

	kafkaGo "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func StartEmailConsumer() {

	reader := kafkaGo.NewReader(
		kafkaGo.ReaderConfig{
			Brokers: []string{
				os.Getenv("KAFKA_BROKER"),
			},

			Topic: kafka.EmailTopic,

			GroupID: "email-worker-group",
		},
	)

	logger.Log.Info(
		"email consumer started",
	)

	for {

		message, err :=
			reader.ReadMessage(
				context.Background(),
			)

		if err != nil {

			logger.Log.Error(
				"failed reading email message",

				zap.Error(err),
			)

			continue
		}

		var payload kafka.EmailPayload

		err = json.Unmarshal(
			message.Value,
			&payload,
		)

		if err != nil {

			logger.Log.Error(
				"failed parsing email payload",

				zap.Error(err),
			)

			continue
		}

		err = email.SendEmail(
			payload.To,
			payload.Subject,
			payload.Body,
		)

		if err != nil {

			logger.Log.Error(
				"failed sending email",

				zap.Error(err),
			)
		}
	}
}
