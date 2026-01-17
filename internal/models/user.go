package models

import "time"

type User struct {
	UserID       int64     `gorm:"primaryKey"`
	Username     string    `gorm:"size:100"`
	LanguagePref string    `gorm:"size:50"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}