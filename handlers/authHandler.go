package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Safiramdhn/project-app-inventaris-golang-safira/models"
	"github.com/Safiramdhn/project-app-inventaris-golang-safira/services"
	"github.com/Safiramdhn/project-app-inventaris-golang-safira/utils"
)

type AuthHandler struct {
	AuthService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

var JsonResp = &utils.JSONResponse{}

func (ah *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		JsonResp.SendError(w, http.StatusMethodNotAllowed, "Method not allowed", r.Method)
		return
	}

	userDTO := models.UserDTO{}
	if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
		JsonResp.SendError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	user, err := ah.AuthService.RegisterUser(&userDTO)
	if err != nil {
		JsonResp.SendError(w, http.StatusInternalServerError, "Failed to register user", err.Error())
		return
	}
	JsonResp.SendCreated(w, user, "User registered")
}

func (ah *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		JsonResp.SendError(w, http.StatusMethodNotAllowed, "Method not allowed", r.Method)
		return
	}

	loginRequest := models.LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		JsonResp.SendError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	token, err := ah.AuthService.LoginUser(&loginRequest)
	if err != nil {
		JsonResp.SendError(w, http.StatusUnauthorized, "Invalid credentials", err.Error())
		return
	}
	JsonResp.SendSuccess(w, token, "User logged in")
}
