package services

import (
	"github.com/kida21/telegram-langbot/internal/models"
	"github.com/kida21/telegram-langbot/internal/repositories"
)

type QuizService struct {
	repo *repositories.QuizRepository
}

func NewQuizService(repo *repositories.QuizRepository) *QuizService {
	return &QuizService{repo: repo}
}

func (s *QuizService) GetRandomQuiz() (*models.Quiz, error) {
	return s.repo.GetRandom()
}