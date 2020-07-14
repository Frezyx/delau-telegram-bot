package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api"
)

func listenAuth(update tgbot.Update, regRequests map[int64]AuthRequest, bot *tgbot.BotAPI, authURL string, checkEmailURL string) tgbot.MessageConfig {
	msg := tgbot.NewMessage(update.Message.Chat.ID, "😿 Я не знаю как ответить на этот вопрос\n"+"Попробуйте что-то отсюда ➡️ /help")

	if thisReq, ok := regRequests[update.Message.Chat.ID]; ok {
		thisReq.Password = update.Message.Text
		regRequests[update.Message.Chat.ID] = thisReq
		if thisReq.Password != "" && thisReq.Email != "" {
			crossMsg := tgbot.NewMessage(
				update.Message.Chat.ID,
				"✅ Вы добавили ваш e-mail.\n"+
					"✅ Вы добавили ваш e-mail и пароль.\n"+
					"Проверяем ваш пароль на сервере.")
			bot.Send(crossMsg)
			buf, err := json.Marshal(regRequests[update.Message.Chat.ID])
			if err != nil {
				msg = tgbot.NewMessage(update.Message.Chat.ID, "⚠️ Произошла ошибка при авторизации")
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
				msg = tgbot.NewMessage(
					update.Message.Chat.ID,
					"⚠️ Вы ввели неверный пароль\n"+
						"Проверьте его правильность и повторите попытку")
			}
		} else {
		}
	} else {
		if update.Message.Text != "" {
			data, _ := http.Get(checkEmailURL + update.Message.Text)
			if data.StatusCode == 200 {
				msg = tgbot.NewMessage(update.Message.Chat.ID, "✅ Вы добавили ваш e-mail.\nТеперь введите ваш пароль :")
				regRequests[update.Message.Chat.ID] = AuthRequest{
					Email:    update.Message.Text,
					Password: "",
					ChatID:   update.Message.Chat.ID,
				}
			} else {
				msg = tgbot.NewMessage(update.Message.Chat.ID, "⚠️ Введенный email не принадлежит аккаунту.")
			}
		}
	}
	return msg
}

func getHelpMessage(update tgbot.Update) tgbot.MessageConfig {
	msg := tgbot.NewMessage(
		update.Message.Chat.ID,
		"🤖 *Команды бота:*\n"+
			"\n"+
			"/login ➡️ Авторизироваться, дла начала работы с ботом\n"+
			"/logout ➡️ Отвязать данный чат от аккаунта и прекратить получать уведомления\n"+
			"/help ➡️ Все запросы бота\n"+
			"/info ➡️ Информация о боте\n")
	msg.ParseMode = "Markdown"
	return msg
}

func getInfoMessage(update tgbot.Update) tgbot.MessageConfig {
	return tgbot.NewMessage(
		update.Message.Chat.ID,
		"Мы часто проводим время в социальных сетях и мессендерах и забываем про личные дела и задачи Delau - проект, созданный для того, чтоб вы смогли получать уведомления о задачах и делах в вашей любимой социальной сети или мессенджере")
}

func getEmailListenerMessage(update tgbot.Update, bot *tgbot.BotAPI) tgbot.MessageConfig {
	msg := tgbot.NewMessage(
		update.Message.Chat.ID,
		"😃 Чтобы получать уведомления о задачах из приложения delau, нам нужно понять кто вы такой.")
	bot.Send(msg)
	return tgbot.NewMessage(update.Message.Chat.ID, "Введите ваш email:")

}

func getStartMessage(update tgbot.Update) tgbot.MessageConfig {
	msg := tgbot.NewMessage(
		update.Message.Chat.ID,
		"*Вас приветствует бот Delau* 😃🖐\n"+
			"\n"+
			"/login ➡️ Чтоб получать уведомления о задачах\n"+
			"/help ➡️ Все запросы бота\n"+
			"/info ➡️ Информация о боте\n")
	msg.ParseMode = "Markdown"
	return msg
}
