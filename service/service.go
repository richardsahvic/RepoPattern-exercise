package service

import (
	"repo"
)

// UserService will be implemented in user_service
type UserService interface {
	Login(username string, password string, role int) (string, error)
	Register(userRegister repo.User, role int) (bool, error)
	ViewProfile(email string, token string) (repo.User, error)
}

var User UserService

func NewService(userRepo repo.UserRepository) {
	User = NewUserService(userRepo)
}
