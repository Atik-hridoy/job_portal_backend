package email

import "log"

// ConsoleSender logs OTP codes instead of sending real emails.
type ConsoleSender struct{}

// NewConsoleSender creates a logging-only email sender.
func NewConsoleSender() *ConsoleSender {
	return &ConsoleSender{}
}

// SendOTP prints the OTP to the application logs so developers can test without SMTP.
func (c *ConsoleSender) SendOTP(to, otp string) error {
	log.Printf("[OTP] To: %s Code: %s", to, otp)
	return nil
}
