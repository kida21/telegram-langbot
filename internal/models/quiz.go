package models

type Quiz struct {
	QuizID        uint   `gorm:"primaryKey;autoIncrement"`
	VocabID       uint   `gorm:"not null"`
	Question      string `gorm:"type:text;not null"`
	Options       string `gorm:"type:jsonb;not null"`
	CorrectOption string `gorm:"size:100;not null"`
}