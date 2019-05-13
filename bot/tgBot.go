package tgBot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/vladkampov/url-shortener-telegram-bot/domain"
	"github.com/vladkampov/url-shortener-telegram-bot/helpers"
	"os"
)

func handleUpdates(bot *tgbotapi.BotAPI, u tgbotapi.UpdateConfig) {
	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		log.Panic(err)
	}


	for update := range updates {
		if update.Message == nil || !helpers.IsUrl(update.Message.Text) {
			// ignore any non-Message Updates and non-url messages
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		shortenedURL := domain.SendUrl(update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, shortenedURL)
		msg.ReplyToMessageID = update.Message.MessageID

		_, err= bot.Send(msg)

		if err != nil {
			log.Panic(err)
		}
	}
}

func Init() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("SHORTENER_BOT_API_TOKEN"))

	if err != nil {
		log.Panic(err)
	}

	if os.Getenv("SHORTENER_BOT_DEBUG") == "true" {
		bot.Debug = true
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	handleUpdates(bot, u)
}