package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/vladkampov/url-shorterer-telegram-bot/bot"
)

func main() {
	log.Printf("We are about to go...")

	tgBot.Init()

	//port := os.Getenv("NOTIFIER_BOT_PORT")
	//
	//if len(port) == 0 {
	//	port = "8080"
	//}
	//
	//log.Printf("Service has started at http://localhost:%s", port)
}