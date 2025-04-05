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

type CampaignHandler struct {
	campaignService *services.CampaignService
	auth            *middleware.Auth
}

func NewCampaignHandler(campaignService *services.CampaignService, auth *middleware.Auth) *CampaignHandler {
	return &CampaignHandler{
		campaignService: campaignService,
		auth:            auth,
	}
}

func (h *CampaignHandler) CreateCampaign(w http.ResponseWriter, r *http.Request) {
	var campaign models.Campaign
	if err := json.NewDecoder(r.Body).Decode(&campaign); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	newCampaign, err := h.campaignService.CreateCampaign(&campaign)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "failed to create campaign")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, newCampaign)
}

func (h *CampaignHandler) GetAllCampaigns(w http.ResponseWriter, r *http.Request) {
	campaigns, err := h.campaignService.GetAllCampaigns()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "failed to fetch campaigns")
		return
	}

	utils.RespondJSON(w, http.StatusOK, campaigns)
}

func (h *CampaignHandler) GetCampaignByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid campaign ID")
		return
	}

	campaign, err := h.campaignService.GetCampaignByID(uint(id))
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "campaign not found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, campaign)
}

func (h *CampaignHandler) DeleteCampaign(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid campaign ID")
		return
	}

	err = h.campaignService.DeleteCampaign(uint(id))
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "failed to delete campaign")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "campaign deleted successfully"})
}
