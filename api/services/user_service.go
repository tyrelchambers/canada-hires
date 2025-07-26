package services

import (
	"canada-hires/models"
	"canada-hires/repos"
	"fmt"
)

type UserService interface {
	GetUserByID(id string) (*models.User, error)
}

type userService struct {
	userRepo repos.UserRepository
}

func NewUserService(userRepo repos.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetUserByID(id string) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}
