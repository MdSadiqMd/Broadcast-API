package models

import (
	"time"

	"gorm.io/gorm"
)

type Campaign struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"uniqueIndex;size:255" json:"name"`
	Role      string         `gorm:"size:50;default:user" json:"role"`
	Contacts  []*Contact     `gorm:"many2many:campaign_audiences;" json:"contacts"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type CampaignAudience struct {
	CampaignID uint `gorm:"primaryKey" json:"campaign_id"`
	ContactID  uint `gorm:"primaryKey" json:"contact_id"`
}
