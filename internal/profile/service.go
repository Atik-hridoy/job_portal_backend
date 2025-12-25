package profile

import (
	"context"
	"errors"
	"strings"

	"job_portal/internal/models"
	"job_portal/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserNotVerified = errors.New("user not verified")
	ErrInvalidGender   = errors.New("invalid gender")
)

// Service coordinates profile operations between repositories and validation.
type Service struct {
	profiles *Repository
	users    *repository.UserRepository
}

func NewService(profiles *Repository, users *repository.UserRepository) *Service {
	return &Service{profiles: profiles, users: users}
}

// CompleteInput represents data supplied by the app when completing a profile.
type CompleteInput struct {
	Email  string
	Name   string
	Gender string
	Phone  string
	Bio    string
	Image  string
}

func (s *Service) CompleteProfile(ctx context.Context, input CompleteInput) (*Profile, error) {
	if s.users == nil {
		return nil, errors.New("user repository not configured")
	}

	user, err := s.users.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	if !user.IsVerified {
		return nil, ErrUserNotVerified
	}

	gender := strings.ToLower(strings.TrimSpace(input.Gender))
	switch gender {
	case "male", "female", "other", "prefer_not_to_say", "":
		// allow empty string if client chooses not to disclose
	default:
		return nil, ErrInvalidGender
	}

	if s.profiles == nil {
		return nil, errors.New("profile repository not configured")
	}

	profileDoc := &Profile{
		Name:   input.Name,
		Gender: gender,
		Phone:  input.Phone,
		Bio:    input.Bio,
		Image:  input.Image,
	}

	profile, err := s.profiles.UpsertByUserID(ctx, user.ID, profileDoc)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

// ViewProfile returns a completed profile along with basic user info.
type ViewProfile struct {
	Profile *Profile     `json:"profile"`
	User    *models.User `json:"user"`
}

func (s *Service) ViewProfile(ctx context.Context, userID primitive.ObjectID) (*ViewProfile, error) {
	if s.users == nil {
		return nil, errors.New("user repository not configured")
	}

	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if s.profiles == nil {
		return nil, errors.New("profile repository not configured")
	}

	profile, err := s.profiles.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &ViewProfile{Profile: profile, User: user}, nil
}
