package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MdSadiqMd/Broadcast-API/internal/api/middleware"
	"github.com/MdSadiqMd/Broadcast-API/internal/services"
	"github.com/MdSadiqMd/Broadcast-API/pkg/utils"
)

type AuthHandler struct {
	userService *services.UserService
	auth        *middleware.Auth
}

func NewAuthHandler(userService *services.UserService, auth *middleware.Auth) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		auth:        auth,
	}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Role     string `json:"role"`
	} `json:"user"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	user, err := h.userService.Authenticate(req.Username, req.Password)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, err := h.auth.GenerateToken(user)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "error generating token")
		return
	}

	resp := loginResponse{
		Token: token,
		User: struct {
			ID       uint   `json:"id"`
			Username string `json:"username"`
			Role     string `json:"role"`
		}{
			ID:       user.ID,
			Username: user.Username,
			Role:     user.Role,
		},
	}

	utils.RespondJSON(w, http.StatusOK, resp)
}

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	user, err := h.userService.CreateUser(req.Username, req.Password, req.Email, "user")
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.auth.GenerateToken(user)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "error generating token")
		return
	}

	resp := loginResponse{
		Token: token,
		User: struct {
			ID       uint   `json:"id"`
			Username string `json:"username"`
			Role     string `json:"role"`
		}{
			ID:       user.ID,
			Username: user.Username,
			Role:     user.Role,
		},
	}

	utils.RespondJSON(w, http.StatusCreated, resp)
}
