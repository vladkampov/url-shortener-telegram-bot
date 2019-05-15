package main

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	tgBot "github.com/vladkampov/url-shortener-telegram-bot/bot"
	"github.com/vladkampov/url-shortener-telegram-bot/domain"
)

func main() {
	log.Printf("We are about to go...")

	err := godotenv.Load()
	if err != nil {
		log.Warn("Error loading .env file")
	}

	_, err = domain.RunDomainGrpcSession()
	if err != nil {
		log.Warnf("Can't start GRPC session: %s", err)
	}

	tgBot.Init()
}
