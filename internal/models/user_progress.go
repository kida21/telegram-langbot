package models

import "time"

type UserProgress struct {
	ProgressID   uint      `gorm:"primaryKey;autoIncrement"`
	UserID       int64     `gorm:"not null"`
	VocabID      uint      `gorm:"not null"`
	Attempts     int       `gorm:"default:0"`
	CorrectCount int       `gorm:"default:0"`
	LastPracticed time.Time `gorm:"autoUpdateTime"`
}