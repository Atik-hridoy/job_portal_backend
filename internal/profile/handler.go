package profile

import (
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

// Handler exposes HTTP endpoints for profile operations.
type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type completeProfileRequest struct {
	Name   string `json:"name"`
	Gender string `json:"gender"`
	Phone  string `json:"phone"`
	Bio    string `json:"bio"`
	Image  string `json:"image"`
}

type completeProfileResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Profile *Profile `json:"profile,omitempty"`
}

// Complete handles POST /api/v1/profile completion submissions from the app.
func (h *Handler) Complete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email, _ := r.Context().Value("user_email").(string)
	if email == "" {
		http.Error(w, "authenticated email not found in context", http.StatusUnauthorized)
		return
	}

	// Parse multipart form (max 32MB)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "failed to parse form data", http.StatusBadRequest)
		return
	}

	// Extract form fields
	payload := completeProfileRequest{
		Name:   r.FormValue("name"),
		Gender: r.FormValue("gender"),
		Phone:  r.FormValue("phone"),
		Bio:    r.FormValue("bio"),
	}

	// Handle image upload
	file, header, err := r.FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		http.Error(w, "failed to retrieve image file", http.StatusBadRequest)
		return
	}
	if file != nil {
		defer file.Close()
		imagePath, err := h.saveImage(file, header)
		if err != nil {
			http.Error(w, "failed to save image", http.StatusInternalServerError)
			return
		}
		payload.Image = imagePath
	}

	profile, err := h.service.CompleteProfile(r.Context(), CompleteInput{
		Email:  email,
		Name:   payload.Name,
		Gender: payload.Gender,
		Phone:  payload.Phone,
		Bio:    payload.Bio,
		Image:  payload.Image,
	})
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(completeProfileResponse{
		Success: true,
		Message: "Profile updated successfully.",
		Profile: profile,
	})
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	switch err {
	case ErrUserNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)
	case ErrUserNotVerified, ErrInvalidGender:
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, "profile not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to complete profile", http.StatusInternalServerError)
	}
}

// saveImage stores uploaded file to uploads directory and returns relative path
func (h *Handler) saveImage(file multipart.File, header *multipart.FileHeader) (string, error) {
	// Create uploads directory if it doesn't exist
	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", err
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	filename := header.Filename
	if filename == "" {
		filename = "profile" + ext
	}

	// Sanitize filename and add timestamp to avoid conflicts
	filename = strings.ReplaceAll(filename, " ", "_")
	timestamp := header.Filename
	if timestamp == "" {
		timestamp = "profile"
	}
	dstPath := filepath.Join(uploadDir, filename)

	// Create destination file
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy uploaded file to destination
	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return "/" + strings.ReplaceAll(dstPath, "\\", "/"), nil
}
