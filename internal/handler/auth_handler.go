package handler

import (
	"encoding/json"
	"net/http"

	"job_portal/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type signUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Phone    int64  `json:"phone"`
}

type signUpResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	Phone int64  `json:"phone"`
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var payload signUpRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	user, err := h.authService.SignUp(r.Context(), service.SignUpInput{
		Email:    payload.Email,
		Password: payload.Password,
		Role:     payload.Role,
		Phone:    payload.Phone,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(signUpResponse{
		ID:    user.ID.Hex(),
		Email: user.Email,
		Role:  user.Role,
		Phone: user.Phone,
	})
}
