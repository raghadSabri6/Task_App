package email

import (
	"bytes"
	"html/template"
	"path/filepath"

	"gopkg.in/gomail.v2"
)

// EmailService handles sending emails
type EmailService struct {
	smtpHost string
	smtpPort int
	smtpUser string
	smtpPass string
	smtpFrom string
}

// NewEmailService creates a new email service
func NewEmailService(smtpHost string, smtpPort int, smtpUser, smtpPass, smtpFrom string) *EmailService {
	return &EmailService{
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		smtpUser: smtpUser,
		smtpPass: smtpPass,
		smtpFrom: smtpFrom,
	}
}

// SendEmail sends an email
func (s *EmailService) SendEmail(to, subject, body string) error {
	// Create message
	m := gomail.NewMessage()
	m.SetHeader("From", s.smtpFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	
	// Create dialer
	d := gomail.NewDialer(s.smtpHost, s.smtpPort, s.smtpUser, s.smtpPass)
	
	// Send email
	return d.DialAndSend(m)
}

// SendTemplateEmail sends an email using a template
func (s *EmailService) SendTemplateEmail(to, subject, templatePath string, data interface{}) error {
	// Parse template
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}
	
	// Execute template
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}
	
	// Send email
	return s.SendEmail(to, subject, body.String())
}

// SendRegistrationEmail sends a registration email
func (s *EmailService) SendRegistrationEmail(to, name string) error {
	// Get template path
	templatePath := filepath.Join("templates", "registration_email_template.html")
	
	// Send email
	return s.SendTemplateEmail(to, "Welcome to Task App", templatePath, map[string]string{
		"Name": name,
	})
}