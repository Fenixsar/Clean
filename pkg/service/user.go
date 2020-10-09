package service

import (
	"gitlab.q123123.net/ligmar/console-back/pkg/repository"
)

type User struct {
	repository.UserRepository
}

func NewUser(repo repository.UserRepository) *User {
	return &User{
		UserRepository: repo,
	}
}
