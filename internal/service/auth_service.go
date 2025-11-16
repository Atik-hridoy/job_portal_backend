package service

import (
	"context"

	"job_portal/internal/models"
	"job_portal/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	users *repository.UserRepository
}

func NewAuthService(users *repository.UserRepository) *AuthService {
	return &AuthService{users: users}
}

type SignUpInput struct {
	Email    string
	Password string
	Role     string
	Phone    int64
}

func (s *AuthService) SignUp(ctx context.Context, input SignUpInput) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:       primitive.NewObjectID(),
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     input.Role,
		Phone:    input.Phone,
	}

	if err := s.users.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
