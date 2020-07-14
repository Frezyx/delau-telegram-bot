package main

import (
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api"
)

var кeyboard = tgbot.NewReplyKeyboard(
	tgbot.NewKeyboardButtonRow(
		tgbot.NewKeyboardButton("/login"),
		tgbot.NewKeyboardButton("/logout"),
	),
	tgbot.NewKeyboardButtonRow(
		tgbot.NewKeyboardButton("/help"),
		tgbot.NewKeyboardButton("/info"),
	),
)
