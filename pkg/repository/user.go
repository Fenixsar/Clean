package repository

import (
	"gitlab.q123123.net/ligmar/boot"
	"gitlab.q123123.net/ligmar/console-back/models"
)

type User struct{}

func NewUserRepository() *User {
	return &User{}
}

func (repo User) GetByTelegramID(id int) (user models.User, err error) {
	err = boot.FindOne(boot.UsersConsoleCollection, id, &user)
	return
}
func (repo User) Store(user models.User) (err error) {
	err = boot.UpsertOne(boot.UsersConsoleCollection, user.TelegramID, user)
	return
}
