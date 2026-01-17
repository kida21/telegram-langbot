package models

import "time"

type Achievement struct {
	AchievementID uint      `gorm:"primaryKey;autoIncrement"`
	UserID        int64     `gorm:"not null"`
	Type          string    `gorm:"size:100;not null"`
	EarnedAt      time.Time `gorm:"autoCreateTime"`
}