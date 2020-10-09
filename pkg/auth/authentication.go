package auth

import (
	"gitlab.q123123.net/ligmar/console-back/models"
	"gitlab.q123123.net/ligmar/console-back/pkg/repository"
)

type Authentication struct {
	repository.UserRepository
}

func NewAuthentication(repo repository.UserRepository) *Authentication {
	return &Authentication{
		repo,
	}
}

func (a *Authentication) Check(id int) (user models.User, err error) {
	user, err = a.GetByTelegramID(id)
	return
}
