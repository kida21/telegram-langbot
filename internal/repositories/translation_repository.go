package repositories

import (
	"github.com/kida21/telegram-langbot/internal/models"
	"gorm.io/gorm"
)

type TranslationRepository struct {
    db *gorm.DB
}

func NewTranslationRepository(db *gorm.DB) *TranslationRepository {
    return &TranslationRepository{db: db}
}

func (r *TranslationRepository) LogTranslation(userID int64, sourceText, targetLang, translatedText string) error {
    entry := models.Translation{
        UserID:         userID,
        SourceText:     sourceText,
        TargetLang:     targetLang,
        TranslatedText: translatedText,
    }
    return r.db.Create(&entry).Error
}

func (r *TranslationRepository) GetHistory(userID int64, limit int) ([]models.Translation, error) {
    var history []models.Translation
    err := r.db.Where("user_id = ?", userID).
        Order("created_at desc").
        Limit(limit).
        Find(&history).Error
    return history, err
}