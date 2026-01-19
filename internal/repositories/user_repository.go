package repositories

import (
	"github.com/kida21/telegram-langbot/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByTelegramID(tgID int64) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, "telegram_id = ?", tgID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) UpdateLanguage(tgID int64, newLang string) error {
    return r.db.Model(&models.User{}).
        Where("telegram_id = ?", tgID).
        Update("language_pref", newLang).Error
}
