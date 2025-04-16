package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	EmailJobStatusQueued   = "queued"
	EmailJobStatusSending  = "sending"
	EmailJobStatusSent     = "sent"
	EmailJobStatusFailed   = "failed"
	EmailJobStatusBounced  = "bounced"
	EmailJobStatusRejected = "rejected"
	EmailJobStatusOpened   = "opened"
	EmailJobStatusClicked  = "clicked"
)

type EmailJob struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	CampaignID    uint           `json:"campaign_id"`
	Campaign      Campaign       `gorm:"foreignkey:CampaignID" json:"campaign"`
	SubscriberID  uint           `json:"subscriber_id"`
	Subscriber    Subscriber     `gorm:"foreignkey:SubscriberID" json:"subscriber"`
	Status        string         `gorm:"size:50;default:queued" json:"status"`
	StatusMessage string         `gorm:"size:255" json:"status_message"`
	Attempts      int            `gorm:"default:0" json:"attempts"`
	SentAt        *time.Time     `json:"sent_at"`
	OpenedAt      *time.Time     `json:"opened_at"`
	ClickedAt     *time.Time     `json:"clicked_at"`
	MessageID     string         `gorm:"size:255" json:"message_id"`
}
