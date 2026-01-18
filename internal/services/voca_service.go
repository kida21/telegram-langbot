package services

import (
	"github.com/kida21/telegram-langbot/internal/models"
	"github.com/kida21/telegram-langbot/internal/repositories"
)

type VocabularyService struct {
	repo *repositories.VocabularyRepository
}

func NewVocabularyService(repo *repositories.VocabularyRepository) *VocabularyService {
	return &VocabularyService{repo: repo}
}

func (s *VocabularyService) GetWordOfTheDay() (*models.Vocabulary, error) {
	return s.repo.GetRandom()
}