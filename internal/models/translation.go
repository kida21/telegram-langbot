package models

import "time"

type Translation struct {
    ID             uint      `gorm:"primaryKey"`
    UserID         int64     `gorm:"index"`
    SourceText     string    `gorm:"type:text;not null"`
    TargetLang     string    `gorm:"size:10;not null"`
    TranslatedText string    `gorm:"type:text;not null"`
    CreatedAt      time.Time `gorm:"autoCreateTime"`
}