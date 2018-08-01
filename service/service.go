package service

import (
	"repo"
)

// UserService will be implemented in user_service
type UserService interface {
	Login(username string, password string) (bool, error)
	Register(userRegister repo.User, role int) (bool, error)
	ViewProfile(email string) (repo.User, error)
}
