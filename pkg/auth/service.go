package auth

import (
	"gitlab.q123123.net/ligmar/console-back/models"
	"gitlab.q123123.net/ligmar/console-back/pkg/repository"
)

type AuthenticationService interface {
	Check(id int) (models.User, error)
}

type Auth struct {
	Authentication AuthenticationService
}

func NewAuth(repos *repository.Repository) *Auth {
	return &Auth{
		NewAuthentication(repos.User),
	}
}
