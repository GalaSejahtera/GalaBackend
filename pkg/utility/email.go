package utility

import (
	"errors"
	"fmt"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
)

func SendPasswordResetEmail(to, name, resetLink string) error {
	data := struct {
		Name      string
		Email     string
		ResetLink string
	}{
		Name:      name,
		Email:     to,
		ResetLink: resetLink,
	}

	file := filepath.Join(os.Getenv("TEMPLATES_PATH"), "password_reset_email.html")

	body, err := ParseHTMLTemplate(file, data)
	if err != nil {
		return errors.New(err.Error() + "filepath: " + file)
	}

	email := email{
		To:      []string{to},
		Subject: "Reset Your Safe Workout Password",
		Body:    body,
	}

	return email.send()
}

// minimum fields required to form a non-spammy email
type email struct {
	To      []string
	Subject string
	Body    string
}

func (e *email) send() error {
	// using galasejahtera email as sender
	username := os.Getenv("EMAIL_ADDRESS")
	password := os.Getenv("EMAIL_PASSWORD")

	smtpHost := os.Getenv("SMTP_SERVER_HOST")
	smtpPort := os.Getenv("SMTP_SERVER_PORT")
	smtpAddr := smtpHost + ":" + smtpPort

	auth := smtp.PlainAuth("", username, password, smtpHost)

	to := fmt.Sprintf("To: %s\n", strings.Join(e.To, ","))
	subject := fmt.Sprintf("Subject: %s\n", e.Subject)
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte(to + subject + mime + e.Body)

	if err := smtp.SendMail(smtpAddr, auth, username, e.To, msg); err != nil {
		return err
	}
	return nil
}
