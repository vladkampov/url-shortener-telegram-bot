package tgBot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/vladkampov/url-shortener-telegram-bot/domain"
	"os"
	"strconv"
	"strings"
)

func handleUpdates(bot *tgbotapi.BotAPI, u tgbotapi.UpdateConfig) {
	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		log.Panic(err)
	}

	webUrl := os.Getenv("SHORTENER_WEB_DOMAIN")
	if len(webUrl) == 0 {
		webUrl = "kmpv.me"
	}

	for update := range updates {
		if update.InlineQuery != nil {
			var results []interface{}

			shortenedURL, err := domain.SendUrl(update.InlineQuery.Query, update.InlineQuery.From.ID)
			var response string
			if err != nil {
				log.Warn(err)
				response = "Uh Oh! Something went wrong. Try again later."
			} else {
				response = shortenedURL
			}

			msg := tgbotapi.NewInlineQueryResultArticleMarkdown(update.InlineQuery.ID, "Short it!", response)

			results = append(results, msg)

			inlineConfig := tgbotapi.InlineConfig{
				InlineQueryID: update.InlineQuery.ID,
				IsPersonal: true,
				CacheTime: 0,
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
				case "start", "help":
					msg.ParseMode = "html"
					msg.Text = "<b>Hey, I'm KMPV URL Shortener</b>\n\n" +
						"Just drop the link including http:// or https:// in it and I'll return the short one.\n" +
						"Notice that it works even in inline mode.\n\n" +
						"<b>Wondering what commands I do?</b>\n" +
 						"- /urls - see the list of your shortened urls.\n" +
						"- /domain customdomain.com – set your domain to work with shortener\n" +
						"- /domain delete – remove custom domain\n" +
						"- /user - see basic information for your user. Here you could verify that your custom domain set properly"
				// TODO: have to decide the payments scheme – implement payments
				case "domain":
					if len(update.Message.CommandArguments()) == 0 {
						msg.Text = "You haven't provide any domain. Use \"/domain myshortdomain.com\" to set custom domain. Use /help for more details"
						break
					}

					user, err := domain.GetUser(update.Message.From.ID)
					if err != nil {
						log.Warnf("Can't get user object for user %s: %s", update.Message.From.UserName, err)
						msg.Text = "Currently can't remove custom domain for you. Try again in a bit!"
						break
					}
					if len(user.CustomDomain) != 0 {
						webUrl = user.CustomDomain
					}

					upperCasedCmdArgs := strings.ToUpper(update.Message.CommandArguments())
					if upperCasedCmdArgs == "DELETE" || upperCasedCmdArgs == "REMOVE" {

						user, err = domain.SetCustomDomain(update.Message.From.ID, "")
						if err != nil {
							log.Warnf("Can't remove custom domain for user %s: %s", update.Message.From.UserName, err)
							msg.Text = "Currently can't remove custom domain for you. Try again in a bit!"
							break
						}

						msg.Text = "Your custom domain were deleted from your profile. Enjoy!"
						break
					}

					user, err = domain.SetCustomDomain(update.Message.From.ID, update.Message.CommandArguments())
					if err != nil {
						log.Warnf("Can't set custom domain for user %s: %s", update.Message.From.UserName, err)
						msg.Text = "Currently can't set custom domain for you. Try again in a bit!"
						break
					}

					webServerIpAddress := os.Getenv("SHORTENER_WEB_IP_ADDRESS")

					msg.ParseMode = "html"
					msg.Text = "Your custom domain successfully added: <b>" + user.CustomDomain + "</b>\n\n" +
						"To finalize this setup please add <b>A</b> record to your domain with hostname (@ and www or any subdomain) and with value <b>" +
						webServerIpAddress + "</b>\n"
				case "user":
					user, err := domain.GetUser(update.Message.From.ID)
					if err != nil {
						log.Warnf("Can't get user object for user %s: %s", update.Message.From.UserName, err)
						msg.Text = "Currently can't get user info for you. Try again in a bit!"
						break
					}
					if len(user.CustomDomain) != 0 {
						webUrl = user.CustomDomain
					}
					msg.ParseMode = "html"
					customDomain := user.CustomDomain
					if len(user.CustomDomain) == 0 {
						customDomain = "You have no custom domain set. Use /domain to set it. Use /help for more info"
					}

					msg.Text = "<b>Here's your basic user info</b>:\n\n" +
						// TODO: have to decide the payments scheme
						//  "<b>Available tokens: </b>" + strconv.FormatInt(int64(user.Tokens), 10) +
						"<b>Custom domain: </b>" + customDomain
				case "urls":
					urls, err := domain.GetUrls(update.Message.From.ID)

					if err != nil {
						log.Warnf("Can't get urls for user %s: %s", update.Message.From.UserName, err)
						msg.Text = "Currently can't get URL's for the user. Try again in a bit"
						break
					}

					if len(urls.Urls) == 0 {
						msg.Text = "You have no urls yet so let's short one quickly!"
						break
					}

					msg.ParseMode = "html"
					message := "<b>Here we go</b>:\n\n"

					for _, url := range urls.Urls {
						message = message + "🔗 " + webUrl + "/" + url.Hash + "\n- visited: <b>" +
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

		shortenedURL, err := domain.SendUrl(update.Message.Text, update.Message.From.ID)
		if err != nil {
			log.Warnf("Error sending message to domain: %s", err)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Uh Oh! Something went wrong. Try again later.")
		} else {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Here's your minified URL: " + shortenedURL)
		}

		msg.ParseMode = "html"
		_, err = bot.Send(msg)
		if err != nil {
			log.Warn(err)
		}
	}
}

func Init() {
	apiToken := strings.TrimSuffix(os.Getenv("SHORTENER_BOT_API_TOKEN"), "\n")
	bot, err := tgbotapi.NewBotAPI(apiToken)

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
