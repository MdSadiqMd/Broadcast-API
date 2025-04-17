package scheduler

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/MdSadiqMd/Broadcast-API/internal/models"
	"github.com/MdSadiqMd/Broadcast-API/internal/services"
	"github.com/MdSadiqMd/Broadcast-API/internal/workers"
	"github.com/MdSadiqMd/Broadcast-API/pkg/config"
	"github.com/MdSadiqMd/Broadcast-API/pkg/email"
	"github.com/go-co-op/gocron"
	"gorm.io/gorm"
)

type Scheduler struct {
	db          *gorm.DB
	config      config.Config
	cron        *gocron.Scheduler
	workers     []*workers.MailWorker
	smtpClient  *email.SMTPClient
	mailService *services.MailService
	mutex       sync.Mutex
}

func NewScheduler(db *gorm.DB, config config.Config) *Scheduler {
	smtpClient := email.NewSMTPClient(email.SMTPConfig{
		Host:     config.SMTP.Host,
		Port:     config.SMTP.Port,
		Username: config.SMTP.Username,
		Password: config.SMTP.Password,
		FromName: config.SMTP.FromName,
		FromAddr: config.SMTP.FromAddr,
		UseTLS:   true,
	})

	mailService := services.NewMailService(db, smtpClient)

	return &Scheduler{
		db:          db,
		config:      config,
		cron:        gocron.NewScheduler(time.UTC),
		workers:     make([]*workers.MailWorker, 0),
		smtpClient:  smtpClient,
		mailService: mailService,
	}
}

func (s *Scheduler) Start() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, err := s.cron.Every(1).Minute().Do(func() {
		s.processCampaigns()
	})
	if err != nil {
		return err
	}

	_, err = s.cron.Every(15).Minutes().Do(func() {
		s.processBouncedEmails()
	})
	if err != nil {
		return err
	}

	_, err = s.cron.Every(1).Hour().Do(func() {
		s.aggregateStats()
	})
	if err != nil {
		return err
	}

	s.cron.StartAsync()
	s.startWorkers()
	log.Println("Scheduler started successfully")
	return nil
}

func (s *Scheduler) Stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.cron.Stop()

	for _, worker := range s.workers {
		worker.Stop()
	}

	log.Println("Scheduler stopped")
}

func (s *Scheduler) startWorkers() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, worker := range s.workers {
		worker.Stop()
	}

	workerCount := s.config.Queue.WorkerCount
	rateLimit := s.config.Queue.RateLimit
	s.workers = make([]*workers.MailWorker, workerCount)

	for i := 0; i < workerCount; i++ {
		worker := workers.NewMailWorker(s.db, s.smtpClient, i+1, rateLimit/workerCount)
		s.workers[i] = worker
		worker.Start()
	}

	log.Printf("Started %d workers with rate limit of %d emails/minute\n", workerCount, rateLimit)
}

func (s *Scheduler) processCampaigns() {
	log.Println("Processing scheduled campaigns...")

	var statusColumnExists bool
	err := s.db.Raw(
		"SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'campaigns' AND column_name = 'status')",
	).Scan(&statusColumnExists).Error

	if err != nil || !statusColumnExists {
		log.Println("Campaign status tracking not available - skipping campaign processing")
		return
	}

	campaignService := services.NewCampaignService(s.db)
	campaigns, err := campaignService.GetScheduledCampaigns()
	if err != nil {
		log.Printf("Error getting scheduled campaigns: %v\n", err)
		return
	}

	for _, campaign := range campaigns {
		log.Printf("Processing campaign: %s (ID: %d)\n", campaign.Name, campaign.ID)
		err = s.db.Model(&campaign).Update("status", models.CampaignStatusProcessing).Error
		if err != nil {
			log.Printf("Error updating campaign status: %v\n", err)
			continue
		}

		var contacts []*models.Contact
		err = s.db.Model(&campaign).Association("Contacts").Find(&contacts)
		if err != nil {
			log.Printf("Error getting contacts for campaign: %v\n", err)
			s.db.Model(&campaign).Updates(map[string]interface{}{
				"status":         models.CampaignStatusError,
				"status_message": err.Error(),
			})
			continue
		}

		if len(contacts) == 0 {
			log.Printf("No contacts found for campaign %d\n", campaign.ID)
			s.db.Model(&campaign).Updates(map[string]interface{}{
				"status":         models.CampaignStatusCompleted,
				"status_message": "No contacts to send to",
			})
			continue
		}

		batchSize := 1000
		totalContacts := len(contacts)
		jobsCreated := 0

		for i := 0; i < totalContacts; i += batchSize {
			end := i + batchSize
			if end > totalContacts {
				end = totalContacts
			}

			batch := contacts[i:end]
			err = s.createEmailJobs(campaign.ID, batch)
			if err != nil {
				log.Printf("Error creating email jobs: %v\n", err)
				continue
			}

			jobsCreated += len(batch)
		}

		now := time.Now()
		s.db.Model(&campaign).Updates(map[string]interface{}{
			"status":         models.CampaignStatusQueued,
			"queued_at":      now,
			"status_message": fmt.Sprintf("Queued %d emails", jobsCreated),
		})

		log.Printf("Queued %d emails for campaign: %s\n", jobsCreated, campaign.Name)
	}
}

func (s *Scheduler) createEmailJobs(campaignID uint, contacts []*models.Contact) error {
	if len(contacts) == 0 {
		return nil
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		jobs := make([]models.EmailJob, len(contacts))

		for i, contact := range contacts {
			jobs[i] = models.EmailJob{
				CampaignID:   campaignID,
				SubscriberID: contact.ID,
				Status:       models.EmailJobStatusQueued,
				CreatedAt:    now,
				UpdatedAt:    now,
			}
		}

		return tx.Create(&jobs).Error
	})
}

func (s *Scheduler) processBouncedEmails() {
	log.Println("Processing bounced emails...")
	// TODO: should add logic for checking a mailbox via IMAP/POP3 or an API
}

func (s *Scheduler) aggregateStats() {
	log.Println("Aggregating email statistics...")
	var campaigns []models.Campaign
	err := s.db.Where("status IN ?", []string{
		models.CampaignStatusRunning,
		models.CampaignStatusCompleted,
	}).Find(&campaigns).Error

	if err != nil {
		log.Printf("Error getting campaigns for stats: %v\n", err)
		return
	}

	for _, campaign := range campaigns {
		var stats = map[string]int{
			"queued":  0,
			"sending": 0,
			"sent":    0,
			"failed":  0,
		}

		err := s.db.Model(&models.EmailJob{}).
			Where("campaign_id = ?", campaign.ID).
			Select("status, count(*) as count").
			Group("status").
			Scan(&[]struct {
				Status string
				Count  int
			}{}).Error

		if err != nil {
			log.Printf("Error getting stats for campaign %d: %v\n", campaign.ID, err)
			continue
		}

		total := stats["queued"] + stats["sending"] + stats["sent"] + stats["failed"]
		sent := stats["sent"]
		failed := stats["failed"]

		if total > 0 && (sent+failed == total) {
			now := time.Now()
			s.db.Model(&campaign).Updates(map[string]interface{}{
				"status":         models.CampaignStatusCompleted,
				"completed_at":   now,
				"status_message": fmt.Sprintf("Completed: %d sent, %d failed", sent, failed),
			})
			log.Printf("Campaign %d marked as completed\n", campaign.ID)
		}
	}
}
