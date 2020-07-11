package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/BurntSushi/toml"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api"
)

var numericKeyboard = tgbot.NewReplyKeyboard(
	tgbot.NewKeyboardButtonRow(
		tgbot.NewKeyboardButton("/login"),
		tgbot.NewKeyboardButton("/logout"),
	),
	tgbot.NewKeyboardButtonRow(
		tgbot.NewKeyboardButton("/help"),
		tgbot.NewKeyboardButton("/info"),
	),
)

func main() {
	config := NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}
	botToken := config.BotToken
	authURL := config.AuthURL
	checkEmailURL := config.CheckEmailURL
	checkTGAuthURL := config.CheckTGAuthURL
	bot, err := tgbot.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	var regRequests map[int64]AuthRequest
	regRequests = make(map[int64]AuthRequest)

	bot.Debug = true

	u := tgbot.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		msg := tgbot.NewMessage(update.Message.Chat.ID, "Я не знаю как ответить на этот вопрос 😿")

		switch update.Message.Text {
		case "/start":
			msg = tgbot.NewMessage(
				update.Message.Chat.ID,
				"Вас приветствует бот Delau 😃\n"+
					"/login - Чтоб получать уведомления о задачах"+
					"/help - Все запросы бота")
		case "/login":
			chatIDStr := strconv.FormatInt(update.Message.Chat.ID, 10)
			data, _ := http.Get(checkTGAuthURL + chatIDStr)
			if data.StatusCode == 200 {
				msg = tgbot.NewMessage(update.Message.Chat.ID, "✅ Вы уже авторизированы.")
			} else {
				//Delete reg data
				msg = tgbot.NewMessage(
					update.Message.Chat.ID,
					"😺Чтобы получать уведомления о задачах из приложения delau, нам нужно понять кто вы такой.")
				bot.Send(msg)
				msg = tgbot.NewMessage(update.Message.Chat.ID, "Введите ваш email:")
			}
		case "/info":
			msg = tgbot.NewMessage(update.Message.Chat.ID, "Мы часто проводим время в социальных сетях и мессендерах и забываем про личные дела и задачи Delau - проект, созданный для того, чтоб вы смогли получать уведомления о задачах и делах в вашей любимой социальной сети или мессенджере")
		case "/help":
			msg = tgbot.NewMessage(
				update.Message.Chat.ID,
				"🤖*Команды бота:*\n"+
					"/login - Авторизироваться, дла начала работы с ботом\n"+
					"/logout - Отвязать данный чат от аккаунта и прекратить получать уведомления\n"+
					"/help - Вывести все команды чата\n"+
					"/info - Информация о приложении\n")
			msg.ParseMode = "Markdown"
		default:
			if thisReq, ok := regRequests[update.Message.Chat.ID]; ok {
				thisReq.Password = update.Message.Text
				regRequests[update.Message.Chat.ID] = thisReq
				if thisReq.Password != "" && thisReq.Email != "" {
					crossMsg := tgbot.NewMessage(
						update.Message.Chat.ID,
						"✅Вы добавили ваш e-mail.\n"+
							"✅Вы добавили ваш пароль.\n"+
							"Проверяем ваш пароль на сервере.")
					bot.Send(crossMsg)
					buf, err := json.Marshal(regRequests[update.Message.Chat.ID])
					if err != nil {
						msg = tgbot.NewMessage(update.Message.Chat.ID, "Произошла ошибка при авторизации 😿")
					}
					data, err := http.Post(authURL, "application/json", bytes.NewBuffer(buf))
					if data.StatusCode == 200 {
						msg = tgbot.NewMessage(
							update.Message.Chat.ID,
							"😃 Отлично\n"+
								"✅ Вы прошли авторизацию.\n"+
								"⏰ Теперь вы получаете уведомления от приложения Dealu")
						regRequests[update.Message.Chat.ID] = AuthRequest{}
					} else {
						msg = tgbot.NewMessage(update.Message.Chat.ID, "Произошла ошибка при авторизации 😿")
					}
				} else {
				}
			} else {
				if update.Message.Text != "" {
					data, _ := http.Get(checkEmailURL + update.Message.Text)
					if data.StatusCode == 200 {
						msg = tgbot.NewMessage(update.Message.Chat.ID, "✅Вы добавили ваш e-mail.\nТеперь введите ваш пароль :")
						regRequests[update.Message.Chat.ID] = AuthRequest{
							Email:    update.Message.Text,
							Password: "",
							ChatID:   update.Message.Chat.ID,
						}
					} else {
						msg = tgbot.NewMessage(update.Message.Chat.ID, "Такого email не существует 😿")
					}
				}
			}
		}
		msg.ReplyMarkup = numericKeyboard
		bot.Send(msg)
	}
}
