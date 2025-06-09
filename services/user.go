package services

import (
	"washo.com/main/repositories"
)

type UserService struct {
	repo *repositories.UserRepository
}
