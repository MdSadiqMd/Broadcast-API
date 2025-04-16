package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	CampaignStatusDraft      = "draft"
	CampaignStatusScheduled  = "scheduled"
	CampaignStatusProcessing = "processing"
	CampaignStatusQueued     = "queued"
	CampaignStatusRunning    = "running"
	CampaignStatusCompleted  = "completed"
	CampaignStatusPaused     = "paused"
	CampaignStatusCancelled  = "cancelled"
	CampaignStatusError      = "error"
)

type Message struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Name          string         `gorm:"size:255" json:"name"`
	Subject       string         `gorm:"size:255" json:"subject"`
	FromEmail     string         `gorm:"size:255" json:"from_email"`
	FromName      string         `gorm:"size:255" json:"from_name"`
	Body          string         `gorm:"type:text" json:"body"`
	Status        string         `gorm:"size:50;default:draft" json:"status"`
	StatusMessage string         `gorm:"size:255" json:"status_message"`
	ScheduledAt   *time.Time     `json:"scheduled_at"`
	StartedAt     *time.Time     `json:"started_at"`
	CompletedAt   *time.Time     `json:"completed_at"`
	QueuedAt      time.Time      `json:"queued_at"`
	ListIDs       []uint         `gorm:"-" json:"list_ids"`
	Lists         []List         `gorm:"many2many:campaign_lists;" json:"lists"`
	CreatedBy     uint           `json:"created_by"`
	User          User           `gorm:"foreignkey:CreatedBy" json:"user"`
}

type List struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"size:255" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Subscribers []Subscriber   `gorm:"many2many:list_subscribers;" json:"subscribers,omitempty"`
}

type Subscriber struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Email     string         `gorm:"uniqueIndex;size:255" json:"email"`
	Name      string         `gorm:"size:255" json:"name"`
	Status    string         `gorm:"size:50;default:enabled" json:"status"`
	Metadata  string         `gorm:"type:jsonb" json:"metadata"`
	Lists     []List         `gorm:"many2many:list_subscribers;" json:"lists,omitempty"`
}
