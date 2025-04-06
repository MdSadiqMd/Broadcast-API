package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/MdSadiqMd/Broadcast-API/internal/api/middleware"
	"github.com/MdSadiqMd/Broadcast-API/internal/models"
	"github.com/MdSadiqMd/Broadcast-API/internal/services"
	"github.com/MdSadiqMd/Broadcast-API/pkg/utils"
	"github.com/go-chi/chi/v5"
)

type BroadcastHandler struct {
	broadcastService *services.BroadcastService
	auth             *middleware.Auth
}

func NewBroadcastHandler(broadcastService *services.BroadcastService, auth *middleware.Auth) *BroadcastHandler {
	return &BroadcastHandler{
		broadcastService: broadcastService,
		auth:             auth,
	}
}

type CreateBroadcastRequest struct {
	Name        string  `json:"name"`
	AudienceID  uint    `json:"audience_id"`
	CampaignID  uint    `json:"campaign_id"`
	UserID      uint    `json:"user_id"`
	From        string  `json:"from"`
	Subject     string  `json:"subject"`
	ReplyTo     string  `json:"reply_to"`
	HTML        string  `json:"html"`
	Text        string  `json:"text"`
	Status      string  `json:"status"`
	ScheduledAt *string `json:"scheduled_at,omitempty"`
	SentAt      *string `json:"sent_at,omitempty"`
}

func (h *BroadcastHandler) CreateBroadcast(w http.ResponseWriter, r *http.Request) {
	var req CreateBroadcastRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	broadcast := models.Broadcast{
		Name:       req.Name,
		AudienceID: req.AudienceID,
		CampaignID: req.CampaignID,
		UserID:     req.UserID,
		From:       req.From,
		Subject:    req.Subject,
		ReplyTo:    req.ReplyTo,
		HTML:       req.HTML,
		Text:       req.Text,
		Status:     req.Status,
	}

	if req.ScheduledAt != nil && *req.ScheduledAt != "" {
		t, err := time.Parse(time.RFC3339, *req.ScheduledAt)
		if err == nil {
			broadcast.ScheduledAt = &t
		}
	}

	if req.SentAt != nil && *req.SentAt != "" {
		t, err := time.Parse(time.RFC3339, *req.SentAt)
		if err == nil {
			broadcast.SentAt = &t
		}
	}

	newBroadcast, err := h.broadcastService.CreateBroadcast(&broadcast)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "failed to create broadcast")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, newBroadcast)
}

func (h *BroadcastHandler) GetBroadcastByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid broadcast ID")
		return
	}

	broadcast, err := h.broadcastService.GetBroadcastByID(uint(id))
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "broadcast not found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, broadcast)
}

func (h *BroadcastHandler) UpdateBroadcast(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid broadcast ID")
		return
	}

	var broadcast models.Broadcast
	if err := json.NewDecoder(r.Body).Decode(&broadcast); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	updatedBroadcast, err := h.broadcastService.UpdateBroadcast(uint(id), &broadcast)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "failed to update broadcast")
		return
	}

	utils.RespondJSON(w, http.StatusOK, updatedBroadcast)
}

type ScheduledPayload struct {
	ScheduledAt string `json:"scheduled_at"`
}

func (h *BroadcastHandler) SendBroadcast(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid broadcast ID")
		return
	}

	var payload ScheduledPayload
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request payload: "+err.Error())
		return
	}

	_, err = h.broadcastService.SendBroadcast(strconv.Itoa(id), payload.ScheduledAt)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "failed to send broadcast: "+err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "broadcast sent successfully"})
}

func (h *BroadcastHandler) ListBroadcasts(w http.ResponseWriter, r *http.Request) {
	broadcasts, err := h.broadcastService.ListBroadcasts()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "failed to fetch broadcasts")
		return
	}

	utils.RespondJSON(w, http.StatusOK, broadcasts)
}

func (h *BroadcastHandler) DeleteBroadcast(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid broadcast ID")
		return
	}

	err = h.broadcastService.DeleteBroadcast(uint(id))
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "failed to delete broadcast")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "broadcast deleted successfully"})
}
