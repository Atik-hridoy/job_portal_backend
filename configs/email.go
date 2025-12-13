package configs

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

var ErrEmailConfigMissing = errors.New("smtp configuration not provided")

type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func LoadEmailConfig() (*EmailConfig, error) {
	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM")

	if host == "" && portStr == "" && username == "" && password == "" && from == "" {
		return nil, ErrEmailConfigMissing
	}

	if host == "" || portStr == "" || username == "" || password == "" || from == "" {
		return nil, fmt.Errorf("SMTP configuration incomplete: set SMTP_HOST, SMTP_PORT, SMTP_USERNAME, SMTP_PASSWORD, SMTP_FROM")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid SMTP_PORT value: %w", err)
	}

	return &EmailConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
	}, nil
}
