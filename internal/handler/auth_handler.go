package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"job_portal/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func (h *AuthHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var payload verifyOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("[VERIFY-OTP] JSON decode error: %v", err)
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	log.Printf("[VERIFY-OTP] Request: Email=%s, OTP=%s", payload.Email, payload.OTP)

	if err := h.authService.VerifyOTP(r.Context(), payload.Email, payload.OTP); err != nil {
		log.Printf("[VERIFY-OTP] Failed for %s: %v", payload.Email, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(verifyOTPResponse{
		Success: true,
		Message: "Account verified successfully.",
	})
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// signUpRequest captures the signup form fields sent from the frontend.
type signUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Phone    string `json:"phone"`
}

type signUpResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type verifyOTPRequest struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

type verifyOTPResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type resendOTPRequest struct {
	Email string `json:"email"`
}

type resendOTPResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type signInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type signInResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	User    interface{} `json:"user,omitempty"`
	Token   string      `json:"token,omitempty"`
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var payload signUpRequest
	// Parse email, password, role selection, and phone number from the request body.
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	// Forward the signup data to the service layer for OTP issuance.
	if err := h.authService.SignUp(r.Context(), service.SignUpInput{
		Email:    payload.Email,
		Password: payload.Password,
		Role:     payload.Role,
		Phone:    payload.Phone,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Respond with a success message indicating the OTP has been dispatched.
	json.NewEncoder(w).Encode(signUpResponse{
		Success: true,
		Message: "OTP has been sent to your email address.",
	})
}

func (h *AuthHandler) ResendOTP(w http.ResponseWriter, r *http.Request) {
	var payload resendOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("[RESEND-OTP] JSON decode error: %v", err)
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	log.Printf("[RESEND-OTP] Request: Email=%s", payload.Email)

	pending, err := h.authService.GetPendingSignup(r.Context(), payload.Email)
	if err != nil {
		log.Printf("[RESEND-OTP] Failed to get pending for %s: %v", payload.Email, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if pending == nil {
		log.Printf("[RESEND-OTP] No pending signup found for %s", payload.Email)
		http.Error(w, "no pending signup found for this email", http.StatusBadRequest)
		return
	}

	otp, err := h.authService.GenerateNewOTP(r.Context(), payload.Email)
	if err != nil {
		log.Printf("[RESEND-OTP] Failed to generate OTP for %s: %v", payload.Email, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("[RESEND-OTP] New OTP generated for %s: %s", payload.Email, otp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resendOTPResponse{
		Success: true,
		Message: "OTP has been resent to your email.",
	})
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var payload signInRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("[SIGN-IN] JSON decode error: %v", err)
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	log.Printf("[SIGN-IN] Request: Email=%s", payload.Email)

	response, err := h.authService.SignIn(r.Context(), service.SignInInput{
		Email:    payload.Email,
		Password: payload.Password,
	})
	if err != nil {
		log.Printf("[SIGN-IN] Failed for %s: %v", payload.Email, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("[SIGN-IN] Success for %s", payload.Email)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(signInResponse{
		Success: true,
		Message: "Sign in successful.",
		User:    response.User,
		Token:   response.Token,
	})
}
