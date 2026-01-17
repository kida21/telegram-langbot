package models

type Vocabulary struct {
	VocabID     uint   `gorm:"primaryKey;autoIncrement"`
	Word        string `gorm:"size:100;not null"`
	Translation string `gorm:"size:100;not null"`
	Example     string `gorm:"type:text"`
	Difficulty  int    `gorm:"default:1"`
}