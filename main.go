package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"job_portal/configs"
	"job_portal/internal/console"
	"job_portal/internal/email"
	"job_portal/internal/handler"
	"job_portal/internal/middleware"
	"job_portal/internal/profile"
	"job_portal/internal/repository"
	"job_portal/internal/service"

	"github.com/joho/godotenv"
)

func main() {
	log.SetFlags(0)
	ui := console.New("Job Portal API")
	ui.Banner("Booting", "Initializing services and establishing connections...")

	if _, err := os.Stat(".env"); err == nil {
		if loadErr := godotenv.Load(); loadErr != nil {
			ui.Status("config", "Failed to load .env file; falling back to system environment")
		} else {
			ui.Status("config", "Loaded environment variables from .env")
		}
	}

	db, err := configs.ConnectDB()
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to close MongoDB connection: %v", err)
		}
	}()
	ui.Status("database", "Connected to MongoDB")

	userRepo := repository.NewUserRepository(db.DB)
	pendingRepo := repository.NewPendingSignupRepository(db.DB)
	profileRepo := profile.NewRepository(db.DB)

	emailCfg, err := configs.LoadEmailConfig()
	var mailer service.EmailSender
	if err != nil {
		if errors.Is(err, configs.ErrEmailConfigMissing) {
			ui.Status("email", "SMTP env not set; OTP codes will be logged locally")
			mailer = email.NewConsoleSender()
		} else {
			log.Fatalf("failed to load SMTP config: %v", err)
		}
	} else {
		mailer = email.NewSMTPSender(emailCfg.Host, emailCfg.Port, emailCfg.Username, emailCfg.Password, emailCfg.From)
		ui.Status("email", "SMTP transport configured")
	}

	authService := service.NewAuthService(userRepo, pendingRepo, mailer)
	authHandler := handler.NewAuthHandler(authService)
	profileService := profile.NewService(profileRepo, userRepo)
	profileHandler := profile.NewHandler(profileService)

	http.HandleFunc("/api/v1/auth/signup", authHandler.SignUp)
	http.HandleFunc("/api/v1/auth/verify-otp", authHandler.VerifyOTP)
	http.HandleFunc("/api/v1/auth/resend-otp", authHandler.ResendOTP)
	http.HandleFunc("/api/v1/auth/signin", authHandler.SignIn)

	// Profile completion endpoint for authenticated app users
	http.Handle("/api/v1/profile", middleware.RequireAuth(http.HandlerFunc(profileHandler.Complete)))

	http.Handle("/api/v1/hirer-only", middleware.RequireEmployer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Hirer only content"}`))
	})))

	http.Handle("/api/v1/job-seeker-only", middleware.RequireJobSeeker(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Job seeker only content"}`))
	})))

	http.Handle("/api/v1/admin-only", middleware.RequireEmployer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Hirer only content (admin functionality removed)"}`))
	})))

	ui.Banner("Server ready", "Listening on http://192.168.1.105:8080")
	if err := http.ListenAndServe("192.168.1.105:8080", nil); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
