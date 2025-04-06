package models

import (
	"time"

	"gorm.io/gorm"
)

type Broadcast struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `json:"name"`
	AudienceID  uint           `json:"audience_id" gorm:"index"`
	CampaignID  uint           `json:"campaign_id" gorm:"index"`
	UserID      uint           `json:"user_id" gorm:"index"`
	From        string         `json:"from"`
	Subject     string         `json:"subject"`
	ReplyTo     string         `json:"reply_to"`
	HTML        string         `json:"html"`
	Text        string         `json:"text"`
	Status      string         `gorm:"default:draft" json:"status"`
	SentAt      *time.Time     `json:"sent_at"`
	Campaign    *Campaign      `gorm:"foreignKey:CampaignID" json:"campaign"`
	User        *User          `gorm:"foreignKey:UserID" json:"user"`
	ScheduledAt *time.Time     `json:"scheduled_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
