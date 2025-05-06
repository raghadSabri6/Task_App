package helperFunc

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"gopkg.in/gomail.v2"
)

func SendGomail(templatePath string, name string, email string) error {
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return fmt.Errorf("email template not found at %s", templatePath)
	}

	var body bytes.Buffer
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("error parsing template: %v", err)
	}

	err = t.Execute(&body, struct{ Name string }{Name: name})
	if err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	sender := os.Getenv("EMAIL_SENDER")
	password := os.Getenv("EMAIL_PASSWORD")
	if sender == "" || password == "" {
		return fmt.Errorf("email configuration missing: EMAIL_SENDER and EMAIL_PASSWORD must be set")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Welcome to Tasks App!")
	m.SetBody("text/html", body.String())

	imagePath := filepath.Join("images", "tasks-app.jpg")
	if _, err := os.Stat(imagePath); err == nil {
		m.Attach(imagePath)
	}

	d := gomail.NewDialer("smtp.gmail.com", 587, sender, password)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}

	return nil
}

func SendSignupEmail(name, email string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting working directory: %v", err)
	}

	templatePath := filepath.Join(currentDir, "templates", "registration_email_template.html")

	return SendGomail(templatePath, name, email)
}
