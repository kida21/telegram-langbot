package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	TelegramID   int64     `gorm:"uniqueIndex;not null"`
	Username     string    `gorm:"uniqueIndex;size:100;not null"`
	LanguagePref string    `gorm:"size:50;not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}