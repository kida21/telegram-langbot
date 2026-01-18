package repositories

import (
	"github.com/kida21/telegram-langbot/internal/models"
	"gorm.io/gorm"
)

type VocabularyRepository struct {
	db *gorm.DB
}

func NewVocabularyRepository(db *gorm.DB) *VocabularyRepository {
	return &VocabularyRepository{db: db}
}

func (r *VocabularyRepository) GetRandom() (*models.Vocabulary, error) {
	var vocab models.Vocabulary
	result := r.db.Order("RANDOM()").First(&vocab)
	if result.Error != nil {
		return nil, result.Error
	}
	return &vocab, nil
}