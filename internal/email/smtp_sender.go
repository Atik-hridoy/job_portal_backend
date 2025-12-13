package email

import (
	"fmt"
	"net/smtp"
)

type SMTPSender struct {
	host     string
	port     int
	username string
	password string
	from     string
}

func NewSMTPSender(host string, port int, username, password, from string) *SMTPSender {
	return &SMTPSender{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (s *SMTPSender) SendOTP(to, otp string) error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: Your verification code\r\n"+
		"MIME-Version: 1.0\r\nContent-Type: text/plain; charset=\"UTF-8\"\r\n\r\n"+
		"Your OTP code is %s. It expires in 5 minutes.\r\n", s.from, to, otp)

	return smtp.SendMail(addr, auth, s.from, []string{to}, []byte(message))
}
