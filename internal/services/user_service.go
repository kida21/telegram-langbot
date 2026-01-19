package services

import (
	"errors"

	"github.com/kida21/telegram-langbot/internal/models"
	"github.com/kida21/telegram-langbot/internal/repositories"
	"gorm.io/gorm"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterOrGet(tgID int64, username string, languagePref string) (*models.User, error) {
	user, err := s.repo.GetByTelegramID(tgID)
	if err == nil {
		return user, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if username == "" {
		return nil, errors.New("username required")
	}

	newUser := &models.User{
		TelegramID:   tgID,
		Username:     username,
		LanguagePref: languagePref,
	}
	if err := s.repo.Create(newUser); err != nil {
		return nil, err
	}
	return newUser, nil
}
