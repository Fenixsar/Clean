package service

import (
	"gitlab.q123123.net/ligmar/console-back/models"
	"gitlab.q123123.net/ligmar/console-back/pkg/repository"
)

type Support struct {
	Repository *repository.Repository
}

func NewSupport(repo *repository.Repository) *Support {
	return &Support{
		repo,
	}
}

func (support Support) GetChatList(skip int) (list models.ChatList, err error) {
	list, err = support.Repository.Support.GetChatList(skip)
	if err != nil {
		return
	}

	var telegramIDs []int64
	for _, l := range list {
		telegramIDs = append(telegramIDs, l.ChatID)
	}

	usersTelegram, err := support.Repository.Telegram.GetByTelegramIDs(telegramIDs)
	if err != nil {
		return
	}

	for i := range list {
		for _, u := range usersTelegram {
			if list[i].ChatID == u.TelegramID {
				list[i].FirstName = u.FirstName
				list[i].LastName = u.LastName
				list[i].Key = u.Key
				break
			}
		}
	}

	return
}
func (support Support) GetMessages(chatID, skip int) (list models.MessageList, err error) {
	list, err = support.Repository.Support.GetMessages(chatID, skip)

	//for _, message := range list {
	//	if message.File != nil {
	//
	//	}
	//}
	//support.telegram.SendMessage(boot.TelegramNotificationID, "Hello")

	return
}
func (support Support) WriteMessage(message models.Message) (err error) {
	err = support.Repository.Support.StoreMessage(message)

	return
}
