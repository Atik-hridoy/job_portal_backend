package main

import (
	"log"
	"net/http"

	"job_portal/configs"
	"job_portal/internal/handler"
	"job_portal/internal/repository"
	"job_portal/internal/service"
)

func main() {
	db, err := configs.ConnectDB()
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to close MongoDB connection: %v", err)
		}
	}()

	userRepo := repository.NewUserRepository(db.DB)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	http.HandleFunc("/api/v1/auth/signup", authHandler.SignUp)

	log.Println("server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
