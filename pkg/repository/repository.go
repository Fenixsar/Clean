package repository

import (
	bootModels "gitlab.q123123.net/ligmar/boot/models"
	"gitlab.q123123.net/ligmar/console-back/models"
)

type SupportRepository interface {
	GetChatList(int) (list models.ChatList, err error)
	GetMessages(chatID, skip int) (list models.MessageList, err error)
	DeleteMessage(messageID int) (err error)
	StoreMessage(message models.Message) (err error)
	MarkAsRead(chatID int) (err error)
}
type TelegramRepository interface {
	GetByTelegramID(id int64) (userTelegram bootModels.UserTelegram, err error)
	GetByTelegramIDs(ids []int64) (userTelegrams []bootModels.UserTelegram, err error)
	Store(userTelegram bootModels.UserTelegram) (err error)
}
type UserRepository interface {
	GetByTelegramID(id int) (user models.User, err error)
	Store(user models.User) (err error)
}

type Repository struct {
	Support  SupportRepository
	Telegram TelegramRepository
	User     UserRepository
}

func NewRepository() *Repository {
	return &Repository{
		NewSupportRepository(),
		NewTelegramRepository(),
		NewUserRepository(),
	}
}
