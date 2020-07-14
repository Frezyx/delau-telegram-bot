package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/BurntSushi/toml"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	config := NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}
	botToken := config.BotToken
	authURL := config.AuthURL
	bot, err := tgbot.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	var regRequests map[int64]AuthRequest
	regRequests = make(map[int64]AuthRequest)

	var isUserStartAuth map[int64]bool
	isUserStartAuth = make(map[int64]bool)

	bot.Debug = true

	u := tgbot.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		msg := tgbot.NewMessage(update.Message.Chat.ID, "Ошибка")

		switch update.Message.Text {
		case "/start":
			msg = getStartMessage(update)
		case "/login":
			chatIDStr := strconv.FormatInt(update.Message.Chat.ID, 10)
			data, _ := http.Get(config.CheckTGAuthURL + chatIDStr)
			if data.StatusCode == 200 {
				msg = tgbot.NewMessage(
					update.Message.Chat.ID,
					"⚠️ Вы уже авторизированы.\n"+
						"Если хотите выйти ➡️ /logout")
				if _, ok := isUserStartAuth[update.Message.Chat.ID]; ok {
					isUserStartAuth[update.Message.Chat.ID] = false
				}
			} else {
				if _, ok := isUserStartAuth[update.Message.Chat.ID]; ok {
					isUserStartAuth[update.Message.Chat.ID] = true
				} else {
					isUserStartAuth[update.Message.Chat.ID] = true
				}
				msg = getEmailListenerMessage(update, bot)
			}
		case "/logout":
			data, _ := http.Get(
				config.LogoutTGURL + strconv.FormatInt(update.Message.Chat.ID, 10))
			if data.StatusCode == 200 {
				msg = tgbot.NewMessage(
					update.Message.Chat.ID,
					"✅ Вы отвязали чат от вашего аккаунта.\n"+
						"Больше вы не получаете уведомления от бота.")
			} else {
				msg = tgbot.NewMessage(
					update.Message.Chat.ID,
					"❌ Ваш чат не привязан ни к отдному аккаунту")
			}
		case "/info":
			msg = getInfoMessage(update)
		case "/help":
			msg = getHelpMessage(update)
		default:
			if _, ok := isUserStartAuth[update.Message.Chat.ID]; ok {
				if isUserStartAuth[update.Message.Chat.ID] {
					msg = listenAuth(update, regRequests, bot, authURL, config.CheckEmailURL)
				}
			} else {
				msg = tgbot.NewMessage(update.Message.Chat.ID, "😿 Я не знаю как ответить на этот вопрос\n"+"Попробуйте что-то отсюда ➡️ /help"+fmt.Sprintf("%t", isUserStartAuth[update.Message.Chat.ID]))
			}
		}
		msg.ReplyMarkup = кeyboard
		bot.Send(msg)
	}
}
