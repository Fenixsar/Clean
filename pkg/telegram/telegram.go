package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type Service interface {
	checkUpdates(updates tgbotapi.UpdatesChannel)
	SendMessage(chatID int64, message string) (tgbotapi.Message, error)
}
