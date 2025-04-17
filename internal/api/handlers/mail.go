package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/MdSadiqMd/Broadcast-API/internal/api/middleware"
	"github.com/MdSadiqMd/Broadcast-API/internal/models"
	"github.com/MdSadiqMd/Broadcast-API/internal/services"
	"github.com/MdSadiqMd/Broadcast-API/pkg/utils"
	"github.com/go-chi/chi/v5"
)

type MailHandler struct {
	mailService *services.MailService
	auth        *middleware.Auth
}

func NewMailHandler(mailService *services.MailService, auth *middleware.Auth) *MailHandler {
	return &MailHandler{
		mailService: mailService,
		auth:        auth,
	}
}

func (h *MailHandler) SendTestEmail(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	if req.Email == "" {
		utils.RespondError(w, http.StatusBadRequest, "email address is required")
		return
	}

	err := h.mailService.SendTestEmail(req.Email)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "failed to send test email: "+err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "test email sent successfully"})
}

func (h *MailHandler) SendTransactionalEmail(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email        string                 `json:"email"`
		Name         string                 `json:"name"`
		Subject      string                 `json:"subject"`
		TemplateName string                 `json:"template_name"`
		Data         map[string]interface{} `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	if req.Email == "" || req.TemplateName == "" || req.Subject == "" {
		utils.RespondError(w, http.StatusBadRequest, "email, template name, and subject are required")
		return
	}

	err := h.mailService.SendTransactionalEmail(req.Email, req.Name, req.Subject, req.TemplateName, req.Data)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "failed to send transactional email: "+err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "transactional email sent successfully"})
}

func (h *MailHandler) ProcessEmailJob(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid job ID")
		return
	}

	job := &models.EmailJob{ID: uint(id)}
	err = h.mailService.ProcessJob(job)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "failed to process email job: "+err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "email job processed successfully"})
}

func (h *MailHandler) ProcessCampaignEmail(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CampaignID uint           `json:"campaign_id"`
		Contact    models.Contact `json:"contact"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	if req.CampaignID == 0 || req.Contact.Email == "" {
		utils.RespondError(w, http.StatusBadRequest, "campaign ID and contact email are required")
		return
	}

	err := h.mailService.ProcessCampaignJob(req.CampaignID, &req.Contact)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "failed to process campaign email: "+err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "campaign email sent successfully"})
}

func (h *MailHandler) BulkSendCampaign(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CampaignID uint             `json:"campaign_id"`
		Contacts   []models.Contact `json:"contacts"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	if req.CampaignID == 0 || len(req.Contacts) == 0 {
		utils.RespondError(w, http.StatusBadRequest, "campaign ID and at least one contact are required")
		return
	}

	successCount := 0
	failureCount := 0

	for _, contact := range req.Contacts {
		err := h.mailService.ProcessCampaignJob(req.CampaignID, &contact)
		if err != nil {
			failureCount++
		} else {
			successCount++
		}
	}

	result := map[string]interface{}{
		"message":       "bulk send complete",
		"success_count": successCount,
		"failure_count": failureCount,
		"total":         len(req.Contacts),
	}

	utils.RespondJSON(w, http.StatusOK, result)
}
