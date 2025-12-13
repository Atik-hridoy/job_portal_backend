package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"job_portal/internal/models"
	"job_portal/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type EmailSender interface {
	SendOTP(to, otp string) error
}

type AuthService struct {
	users       *repository.UserRepository
	pendings    *repository.PendingSignupRepository
	emailSender EmailSender
	otpTTL      time.Duration
	jwtSecret   string
}

type SignInInput struct {
	Email    string
	Password string
}

type SignInResponse struct {
	User  *models.User `json:"user"`
	Token string       `json:"token,omitempty"`
}

func NewAuthService(users *repository.UserRepository, pendings *repository.PendingSignupRepository, sender EmailSender) *AuthService {
	return &AuthService{
		users:       users,
		pendings:    pendings,
		emailSender: sender,
		otpTTL:      5 * time.Minute,
		jwtSecret:   "your-secret-key", // TODO: Move to environment variable
	}
}

type SignUpInput struct {
	Email    string
	Password string
	Role     string
	Phone    string
}

func (s *AuthService) SignUp(ctx context.Context, input SignUpInput) error {
	if s.users != nil {
		existingUser, err := s.users.GetByEmail(ctx, input.Email)
		if err != nil {
			return err
		}
		if existingUser != nil {
			return fmt.Errorf("account already verified for this email")
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	otp, err := generateOTP()
	if err != nil {
		return err
	}

	if s.emailSender == nil {
		return fmt.Errorf("email sender not configured")
	}

	now := primitive.NewDateTimeFromTime(time.Now())
	pending := &models.PendingSignup{
		Email:     input.Email,
		Password:  string(hashedPassword),
		Role:      input.Role,
		Phone:     input.Phone,
		OTP:       otp,
		OTPExpiry: primitive.NewDateTimeFromTime(time.Now().Add(s.otpTTL)),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if s.pendings == nil {
		return fmt.Errorf("pending signup repository not configured")
	}

	if err := s.pendings.Upsert(ctx, pending); err != nil {
		return err
	}

	if err := s.emailSender.SendOTP(pending.Email, otp); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) VerifyOTP(ctx context.Context, email, otp string) error {
	if s.pendings == nil {
		return fmt.Errorf("pending signup repository not configured")
	}

	pending, err := s.pendings.GetByEmail(ctx, email)
	if err != nil {
		return err
	}
	if pending == nil {
		return fmt.Errorf("no pending signup found for this email")
	}

	if pending.OTP != otp {
		return fmt.Errorf("invalid OTP")
	}

	if pending.OTPExpiry != 0 && time.Now().After(pending.OTPExpiry.Time()) {
		return fmt.Errorf("otp expired")
	}

	if s.users == nil {
		return fmt.Errorf("user repository not configured")
	}

	user := &models.User{
		Email:      pending.Email,
		Password:   pending.Password,
		Role:       pending.Role,
		Phone:      pending.Phone,
		IsVerified: true,
		CreatedAt:  primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt:  primitive.NewDateTimeFromTime(time.Now()),
	}

	if err := s.users.CreateUser(ctx, user); err != nil {
		return err
	}

	if err := s.pendings.Delete(ctx, pending.Email); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) GetPendingSignup(ctx context.Context, email string) (*models.PendingSignup, error) {
	if s.pendings == nil {
		return nil, fmt.Errorf("pending signup repository not configured")
	}
	return s.pendings.GetByEmail(ctx, email)
}

func (s *AuthService) GenerateNewOTP(ctx context.Context, email string) (string, error) {
	if s.pendings == nil {
		return "", fmt.Errorf("pending signup repository not configured")
	}

	pending, err := s.pendings.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if pending == nil {
		return "", fmt.Errorf("no pending signup found for this email")
	}

	otp, err := generateOTP()
	if err != nil {
		return "", err
	}

	pending.OTP = otp
	pending.OTPExpiry = primitive.NewDateTimeFromTime(time.Now().Add(s.otpTTL))
	pending.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	if err := s.pendings.Upsert(ctx, pending); err != nil {
		return "", err
	}

	if s.emailSender == nil {
		return "", fmt.Errorf("email sender not configured")
	}

	if err := s.emailSender.SendOTP(pending.Email, otp); err != nil {
		return "", err
	}

	return otp, nil
}

func (s *AuthService) SignIn(ctx context.Context, input SignInInput) (*SignInResponse, error) {
	if s.users == nil {
		return nil, fmt.Errorf("user repository not configured")
	}

	user, err := s.users.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("no verified user found with this email")
	}

	if !user.IsVerified {
		return nil, fmt.Errorf("account not verified. Please verify your email first")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	token, err := s.generateJWT(user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	return &SignInResponse{
		User:  user,
		Token: token,
	}, nil
}

func (s *AuthService) generateJWT(email, role string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func generateOTP() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}
