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

func (h *BroadcastHandler) CreateBroadcast(w http.ResponseWriter, r *http.Request) {
	var broadcast models.Broadcast
	if err := json.NewDecoder(r.Body).Decode(&broadcast); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request payload")
		return
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

func (h *BroadcastHandler) SendBroadcast(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid broadcast ID")
		return
	}

	var scheduled_at string
	err = json.NewDecoder(r.Body).Decode(&scheduled_at)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request payload: "+err.Error())
		return
	}

	_, err = h.broadcastService.SendBroadcast(strconv.Itoa(id), scheduled_at)
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
