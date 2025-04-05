package models

import (
	"time"

	"gorm.io/gorm"
)

type Campaign struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"uniqueIndex;size:255" json:"name"`
	Role      string         `gorm:"size:50;default:user" json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
