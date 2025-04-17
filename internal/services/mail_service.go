package services

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/MdSadiqMd/Broadcast-API/internal/models"
	"github.com/MdSadiqMd/Broadcast-API/pkg/email"
	"gorm.io/gorm"
)

type MailService struct {
	db         *gorm.DB
	smtpClient *email.SMTPClient
	templates  map[string]*template.Template
}

func NewMailService(db *gorm.DB, smtpClient *email.SMTPClient) *MailService {
	return &MailService{
		db:         db,
		smtpClient: smtpClient,
		templates:  make(map[string]*template.Template),
	}
}

func (s *MailService) SendTestEmail(toEmail string) error {
	if toEmail == "" {
		return errors.New("recipient email is required")
	}

	message := email.Message{
		To:      toEmail,
		Subject: "Test Email from Broadcast API",
		HTML:    "<html><body><h1>Test Email</h1><p>This is a test email from your Broadcast API system.</p><p>If you're receiving this, your email configuration is working properly!</p></body></html>",
		Text:    "Test Email\n\nThis is a test email from your Broadcast API system.\nIf you're receiving this, your email configuration is working properly!",
		Headers: map[string]string{"X-Test": "true"},
	}

	_, err := s.smtpClient.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send test email: %w", err)
	}

	return nil
}

func (s *MailService) SendTransactionalEmail(
	toEmail string,
	toName string,
	subject string,
	templateName string,
	data map[string]interface{}) error {

	tmpl, err := s.getTemplate(templateName)
	if err != nil {
		return err
	}

	if data == nil {
		data = make(map[string]interface{})
	}
	data["toEmail"] = toEmail
	data["toName"] = toName

	var htmlBuf bytes.Buffer
	err = tmpl.Execute(&htmlBuf, data)
	if err != nil {
		return fmt.Errorf("template execution error: %w", err)
	}
	textContent := "Please view this email with an HTML-capable email client."

	message := email.Message{
		To:      toEmail,
		Subject: subject,
		HTML:    htmlBuf.String(),
		Text:    textContent,
		Headers: map[string]string{"X-Email-Type": "transactional"},
	}

	_, err = s.smtpClient.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	err = s.logTransactionalEmail(toEmail, subject, templateName)
	if err != nil {
		log.Printf("Failed to log transactional email: %v", err)
	}

	return nil
}

func (s *MailService) ProcessJob(job *models.EmailJob) error {
	err := s.db.Preload("Campaign").Preload("Subscriber").First(job, job.ID).Error
	if err != nil {
		return fmt.Errorf("error loading job data: %w", err)
	}

	var message models.Message
	err = s.db.Where("id = ?", job.Campaign.ID).First(&message).Error
	if err != nil {
		return fmt.Errorf("error loading message data: %w", err)
	}

	data := map[string]interface{}{
		"subscriber": job.Subscriber,
		"campaign":   job.Campaign,
		"message":    message,
		"date":       time.Now().Format("2006-01-02"),
		"unsubscribeURL": fmt.Sprintf("http://example.com/unsubscribe?email=%s&uuid=%d",
			job.Subscriber.Email, job.Subscriber.ID),
	}

	subjectTmpl, err := template.New("subject").Parse(message.Subject)
	if err != nil {
		return fmt.Errorf("subject template parse error: %w", err)
	}
	var subjectBuf bytes.Buffer
	err = subjectTmpl.Execute(&subjectBuf, data)
	if err != nil {
		return fmt.Errorf("subject template execution error: %w", err)
	}
	subject := subjectBuf.String()
	htmlTmpl, err := template.New("html").Parse(message.Body)
	if err != nil {
		return fmt.Errorf("html template parse error: %w", err)
	}
	var htmlBuf bytes.Buffer
	err = htmlTmpl.Execute(&htmlBuf, data)
	if err != nil {
		return fmt.Errorf("html template execution error: %w", err)
	}
	htmlContent := htmlBuf.String()
	textContent := "Please view this email with an HTML-capable email client."

	emailMessage := email.Message{
		FromEmail: message.FromEmail,
		FromName:  message.FromName,
		To:        job.Subscriber.Email,
		Subject:   subject,
		HTML:      htmlContent,
		Text:      textContent,
		Headers: map[string]string{
			"X-Campaign-ID":   fmt.Sprintf("%d", job.CampaignID),
			"X-Subscriber-ID": fmt.Sprintf("%d", job.SubscriberID),
		},
	}

	messageID, err := s.smtpClient.Send(emailMessage)
	if err != nil {
		job.Status = models.EmailJobStatusFailed
		job.StatusMessage = err.Error()
		s.db.Save(job)
		return fmt.Errorf("failed to send email: %w", err)
	}

	now := time.Now()
	job.Status = models.EmailJobStatusSent
	job.SentAt = &now
	job.MessageID = messageID
	err = s.db.Save(job).Error
	if err != nil {
		log.Printf("Failed to update job status: %v", err)
	}

	return nil
}

func (s *MailService) ProcessCampaignJob(campaignID uint, contact *models.Contact) error {
	var campaign models.Campaign
	err := s.db.First(&campaign, campaignID).Error
	if err != nil {
		return fmt.Errorf("error loading campaign data: %w", err)
	}

	var message models.Message
	err = s.db.Where("id = ?", campaignID).First(&message).Error
	if err != nil {
		return fmt.Errorf("error loading message data: %w", err)
	}

	data := map[string]interface{}{
		"contact":  contact,
		"campaign": campaign,
		"message":  message,
		"date":     time.Now().Format("2006-01-02"),
		"unsubscribeURL": fmt.Sprintf("http://example.com/unsubscribe?email=%s&uuid=%d",
			contact.Email, contact.ID),
	}

	subjectTmpl, err := template.New("subject").Parse(message.Subject)
	if err != nil {
		return fmt.Errorf("subject template parse error: %w", err)
	}
	var subjectBuf bytes.Buffer
	err = subjectTmpl.Execute(&subjectBuf, data)
	if err != nil {
		return fmt.Errorf("subject template execution error: %w", err)
	}
	subject := subjectBuf.String()

	htmlTmpl, err := template.New("html").Parse(message.Body)
	if err != nil {
		return fmt.Errorf("html template parse error: %w", err)
	}
	var htmlBuf bytes.Buffer
	err = htmlTmpl.Execute(&htmlBuf, data)
	if err != nil {
		return fmt.Errorf("html template execution error: %w", err)
	}
	htmlContent := htmlBuf.String()

	textContent := "Please view this email with an HTML-capable email client."

	emailMessage := email.Message{
		FromEmail: message.FromEmail,
		FromName:  message.FromName,
		To:        contact.Email,
		Subject:   subject,
		HTML:      htmlContent,
		Text:      textContent,
		Headers: map[string]string{
			"X-Campaign-ID": fmt.Sprintf("%d", campaignID),
			"X-Contact-ID":  fmt.Sprintf("%d", contact.ID),
		},
	}

	messageID, err := s.smtpClient.Send(emailMessage)
	if err != nil {
		log.Printf("Failed to send email to contact %d: %v", contact.ID, err)
		job := &models.EmailJob{
			CampaignID:    campaignID,
			SubscriberID:  contact.ID,
			Status:        models.EmailJobStatusFailed,
			StatusMessage: err.Error(),
			Attempts:      1,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		s.db.Create(job)

		return fmt.Errorf("failed to send email: %w", err)
	}

	now := time.Now()
	job := &models.EmailJob{
		CampaignID:   campaignID,
		SubscriberID: contact.ID,
		Status:       models.EmailJobStatusSent,
		MessageID:    messageID,
		Attempts:     1,
		SentAt:       &now,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err = s.db.Create(job).Error
	if err != nil {
		log.Printf("Failed to create job record: %v", err)
	}

	return nil
}

func (s *MailService) getTemplate(name string) (*template.Template, error) {
	if tmpl, ok := s.templates[name]; ok {
		return tmpl, nil
	}

	var templateContent string
	err := s.db.Model(models.Template{}).
		Where("name = ?", name).
		Select("content").
		Scan(&templateContent).Error
	if err != nil {
		return nil, fmt.Errorf("template not found: %w", err)
	}

	tmpl, err := template.New(name).Parse(templateContent)
	if err != nil {
		return nil, fmt.Errorf("template parse error: %w", err)
	}

	s.templates[name] = tmpl
	return tmpl, nil
}

func (s *MailService) logTransactionalEmail(toEmail, subject, templateName string) error {
	log := models.EmailLog{
		Email:    toEmail,
		Subject:  subject,
		Template: templateName,
		Type:     "transactional",
		SentAt:   time.Now(),
		Status:   "sent",
	}

	return s.db.Create(&log).Error
}
