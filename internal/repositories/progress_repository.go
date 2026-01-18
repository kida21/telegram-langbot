package repositories

import (
	"github.com/kida21/telegram-langbot/internal/models"
	"gorm.io/gorm"
)

type ProgressRepository struct {
	db *gorm.DB
}

func NewProgressRepository(db *gorm.DB) *ProgressRepository {
	return &ProgressRepository{db: db}
}

func (r *ProgressRepository) GetByUser(userID int64) ([]models.UserProgress, error) {
	var progress []models.UserProgress
	result := r.db.Where("user_id = ?", userID).Find(&progress)
	return progress, result.Error
}