package tgBot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/vladkampov/url-shortener-telegram-bot/domain"
	"github.com/vladkampov/url-shortener-telegram-bot/helpers"
	"os"
	"strconv"
)

func handleUpdates(bot *tgbotapi.BotAPI, u tgbotapi.UpdateConfig) {
	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		log.Panic(err)
	}


	for update := range updates {
		if update.InlineQuery != nil {
			var results []interface{}
			query := update.InlineQuery.Query

			var msg tgbotapi.InlineQueryResultArticle
			if !helpers.IsUrl(query) {
				msg = tgbotapi.NewInlineQueryResultArticleMarkdown(update.InlineQuery.ID, "Put the full url here", "Put the full URL here: it has to be with http:// or https://")
			} else {
				shortenedURL, err := domain.SendUrl(update.InlineQuery.Query, update.InlineQuery.From.ID)

				var response string
				if err != nil {
					log.Warn(err)
					response = "Uh Oh! Something went wrong. Try again later."
				} else {
					response = shortenedURL
				}

				msg = tgbotapi.NewInlineQueryResultArticleMarkdown(update.InlineQuery.ID, "Short it!", response)
			}
			results = append(results, msg)

			inlineConfig := tgbotapi.InlineConfig{
				InlineQueryID: update.InlineQuery.ID,
				IsPersonal:    true,
				CacheTime:     0,
				Results: results,
			}
			_, err := bot.AnswerInlineQuery(inlineConfig)
			if err != nil {
				log.Println(err)
				continue
			}
			continue
		}

		if update.Message == nil {
			// ignore any non-Message Updates and non-url messages
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)


		var msg tgbotapi.MessageConfig
		msg.ReplyToMessageID = update.Message.MessageID

		if update.Message.IsCommand() {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
				case "urls":
					urls, err := domain.GetUrls(update.Message.From.ID)

					if err != nil {
						log.Warnf("Can't get urls for user %s: %s", update.Message.From.UserName, err)
						continue
					}

					if len(urls.Urls) == 0 {
						msg.Text = "You have no urls yet so let's short one quickly!"
						break
					}

					msg.ParseMode = "html"
					message := "<b>Here we go</b>:\n\n"

					for _, url := range urls.Urls {
						log.Println(url)
						message = message + "🔗 kmpv.me/" + url.Hash + "\n- visited: <b>" +
							strconv.FormatInt(int64(url.Visited), 10) +
							"</b> time(s) \n- <a href=\"" + url.Url + "\">link</a>\n\n"
					}

					msg.Text = message
				default:
					msg.Text = "I don't know that command yet :)"
			}

			_, err = bot.Send(msg)
			if err != nil {
				log.Warn(err)
			}

			continue
		}

		if !helpers.IsUrl(update.Message.Text) {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Seems like this is not a valid URL. Try again")
		} else {
			shortenedURL, err := domain.SendUrl(update.Message.Text, update.Message.From.ID)

			if err != nil {
				log.Warnf("Error sending message to domain: %s", err)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Uh Oh! Something went wrong. Try again later.")
			} else {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Here's your minified URL: " + shortenedURL)
			}
		}

		_, err = bot.Send(msg)

		if err != nil {
			log.Warn(err)
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