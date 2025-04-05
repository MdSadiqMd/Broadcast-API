package models

import (
	"time"

	"gorm.io/gorm"
)

type Contact struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	Email       string         `gorm:"uniqueIndex;size:255" json:"email"`
	UnSubscribe bool           `json:"unsubscribe"`
	Campaigns   []*Campaign    `gorm:"many2many:campaign_audiences;" json:"campaigns"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
