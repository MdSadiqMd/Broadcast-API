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

type ContactHandler struct {
	contactService *services.ContactService
	auth           *middleware.Auth
}

func NewContactHandler(contactService *services.ContactService, auth *middleware.Auth) *ContactHandler {
	return &ContactHandler{
		contactService: contactService,
		auth:           auth,
	}
}

type GetAllContactsRequest struct {
	AudienceId uint `json:"audienceId"`
}

func (h *ContactHandler) CreateContact(w http.ResponseWriter, r *http.Request) {
	var contact models.Contact
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	newContact, err := h.contactService.CreateContact(&contact)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "failed to create contact")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, newContact)
}

func (h *ContactHandler) GetAllContacts(w http.ResponseWriter, r *http.Request) {
	var request GetAllContactsRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	contacts, err := h.contactService.GetAllContacts(request.AudienceId)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "failed to fetch contacts")
		return
	}

	utils.RespondJSON(w, http.StatusOK, contacts)
}

func (h *ContactHandler) GetContactByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid contact ID")
		return
	}

	contact, err := h.contactService.GetContactByID(uint(id))
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "contact not found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, contact)
}

func (h *ContactHandler) UpdateContact(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid contact ID")
		return
	}

	var contact models.Contact
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	updatedContact, err := h.contactService.UpdateContact(uint(id), &contact)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "failed to update contact")
		return
	}

	utils.RespondJSON(w, http.StatusOK, updatedContact)
}

func (h *ContactHandler) DeleteContact(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid contact ID")
		return
	}

	err = h.contactService.DeleteContact(uint(id))
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "failed to delete contact")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "contact deleted successfully"})
}
