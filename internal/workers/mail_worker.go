package workers

import (
	"log"
	"sync"
	"time"

	"github.com/MdSadiqMd/Broadcast-API/internal/models"
	"github.com/MdSadiqMd/Broadcast-API/internal/services"
	"github.com/MdSadiqMd/Broadcast-API/pkg/email"
	"gorm.io/gorm"
)

type MailWorker struct {
	db          *gorm.DB
	smtpClient  *email.SMTPClient
	mailService *services.MailService
	workerID    int
	wg          *sync.WaitGroup
	stopChan    chan struct{}
	rateLimiter *RateLimiter
	running     bool
	mutex       sync.Mutex
}

func NewMailWorker(db *gorm.DB, smtpClient *email.SMTPClient, workerID int, rateLimit int) *MailWorker {
	mailService := services.NewMailService(db, smtpClient)

	return &MailWorker{
		db:          db,
		smtpClient:  smtpClient,
		mailService: mailService,
		workerID:    workerID,
		wg:          &sync.WaitGroup{},
		stopChan:    make(chan struct{}),
		rateLimiter: NewRateLimiter(rateLimit),
		running:     false,
	}
}

func (w *MailWorker) Start() {
	w.mutex.Lock()
	if w.running {
		w.mutex.Unlock()
		return
	}
	w.running = true
	w.mutex.Unlock()

	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		w.processJobs()
	}()

	log.Printf("Worker %d started\n", w.workerID)
}

func (w *MailWorker) Stop() {
	w.mutex.Lock()
	if !w.running {
		w.mutex.Unlock()
		return
	}
	w.running = false
	close(w.stopChan)
	w.mutex.Unlock()

	log.Printf("Worker %d stopping...\n", w.workerID)
	w.wg.Wait()
	log.Printf("Worker %d stopped\n", w.workerID)
}

func (w *MailWorker) processJobs() {
	for {
		select {
		case <-w.stopChan:
			return
		default:
			job, err := w.getNextJob()
			if err != nil {
				time.Sleep(5 * time.Second)
				continue
			}

			w.rateLimiter.Wait()
			w.processJob(job)
		}
	}
}

func (w *MailWorker) getNextJob() (*models.EmailJob, error) {
	var job models.EmailJob
	err := w.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("status = ?", models.EmailJobStatusQueued).
			Order("created_at asc").
			First(&job)

		if result.Error != nil {
			return result.Error
		}

		job.Status = models.EmailJobStatusSending
		job.Attempts++
		job.UpdatedAt = time.Now()

		return tx.Save(&job).Error
	})

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}

	return &job, nil
}

func (w *MailWorker) processJob(job *models.EmailJob) {
	log.Printf("Worker %d processing job %d\n", w.workerID, job.ID)
	err := w.mailService.ProcessJob(job)

	if err != nil {
		log.Printf("Worker %d error processing job %d: %v\n", w.workerID, job.ID, err)
		w.markJobAsFailed(job, err.Error())
		return
	}

	log.Printf("Worker %d successfully processed job %d\n", w.workerID, job.ID)
}

func (w *MailWorker) markJobAsFailed(job *models.EmailJob, errorMessage string) {
	if job.Attempts >= 3 {
		job.Status = models.EmailJobStatusFailed
	} else {
		job.Status = models.EmailJobStatusQueued
	}

	job.StatusMessage = errorMessage
	job.UpdatedAt = time.Now()

	if err := w.db.Save(job).Error; err != nil {
		log.Printf("Worker %d error updating failed job status: %v\n", w.workerID, err)
	}
}

type RateLimiter struct {
	ticker *time.Ticker
}

func NewRateLimiter(ratePerMinute int) *RateLimiter {
	if ratePerMinute <= 0 {
		ratePerMinute = 60
	}

	interval := time.Minute / time.Duration(ratePerMinute)
	return &RateLimiter{
		ticker: time.NewTicker(interval),
	}
}

func (r *RateLimiter) Wait() {
	<-r.ticker.C
}
