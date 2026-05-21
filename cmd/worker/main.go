package main

import (
	"notification-service/internal/logger"
	"notification-service/internal/worker"

	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load()

	logger.InitLogger()

	defer logger.Log.Sync()

	go worker.StartOTPConsumer()

	go worker.StartWhatsAppConsumer()

	go worker.StartEmailConsumer()

	select {}
}
