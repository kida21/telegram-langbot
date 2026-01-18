package repositories

import (
	"github.com/kida21/telegram-langbot/internal/models"
	"gorm.io/gorm"
)

type QuizRepository struct {
	db *gorm.DB
}

func NewQuizRepository(db *gorm.DB) *QuizRepository {
	return &QuizRepository{db: db}
}

func (r *QuizRepository) GetRandom() (*models.Quiz, error) {
	var quiz models.Quiz
	result := r.db.Order("RANDOM()").First(&quiz)
	if result.Error != nil {
		return nil, result.Error
	}
	return &quiz, nil
}