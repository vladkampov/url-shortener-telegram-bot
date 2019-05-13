package main

import (
	log "github.com/sirupsen/logrus"
	tgBot "github.com/vladkampov/url-shortener-telegram-bot/bot"
	"github.com/vladkampov/url-shortener-telegram-bot/domain"
)

func main() {
	log.Printf("We are about to go...")
	domain.InitDomainGrpcSession()
	tgBot.Init()
}
