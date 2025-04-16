package models

import (
	"time"

	"gorm.io/gorm"
)

type Template struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `gorm:"uniqueIndex;size:255" json:"name"`
	Content   string         `gorm:"type:text" json:"content"`
	Type      string         `gorm:"size:50;default:html" json:"type"`
}

type EmailLog struct {
	ID       uint      `gorm:"primarykey" json:"id"`
	Email    string    `gorm:"size:255;index" json:"email"`
	Subject  string    `gorm:"size:255" json:"subject"`
	Template string    `gorm:"size:255" json:"template"`
	Type     string    `gorm:"size:50" json:"type"`
	SentAt   time.Time `json:"sent_at"`
	Status   string    `gorm:"size:50" json:"status"`
}
