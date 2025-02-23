package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"test-ai-api/stores"
	"test-ai-api/types"
	"test-ai-api/utils"
)

type AuthHandler struct {
	userStore *stores.UserStore
}

func NewAuthHandler(userStore *stores.UserStore) *AuthHandler {
	return &AuthHandler{userStore: userStore}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// Handle preflight request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := h.userStore.Login(loginRequest.Email, loginRequest.Password)
	fmt.Println(user)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// Handle preflight request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var user types.UserRegister
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate required fields
	if user.Email == "" || user.Password == "" || user.FirstName == "" || user.LastName == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	// Check if email already exists
	_, err := h.userStore.GetByEmail(user.Email)
	if err == nil {
		utils.RespondWithError(w, http.StatusConflict, "Email already exists")
		return
	}

	newUser, err := h.userStore.Register(user)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(newUser)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"token": token,
		"user":  newUser,
	})
}

func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int64)
	user, err := h.userStore.GetByID(userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, user)
}
