package helpers

import (
	"bytes"
	"net/smtp"
	"os"
	"strings"
	"text/template"
)

type Request struct {
	to      []string
	subject string
	body    string
	from    string
}

func (r *Request) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		println("parse error", err)
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		println("read error", err)
		return err
	}
	r.body = buf.String()
	return nil
}

func NewRequest(to []string, subject, body string) *Request {
	return &Request{
		to:      to,
		subject: subject,
		body:    body,
		from:    os.Getenv("MAIL_FROM"),
	}
}

func (r *Request) SendEmail() (bool, error) {
	fromEmail := r.from
	fromPass := os.Getenv("SMTP_PASSWORD")
	emailIdentity := os.Getenv("EMAIL_IDENTITY")
	emailClient := os.Getenv("SMTP_CLIENT")
	emailPort := os.Getenv("SMTP_PORT")
	auth := smtp.PlainAuth(emailIdentity, fromEmail, fromPass, emailClient)

	toHeader := strings.Join(r.to, ", ")
	subject := "Subject: " + r.subject + "\n"
	from := "From: " + r.from + "\n"
	to := "To: " + toHeader + "\n"
	// send html email
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	message := []byte(from + to + subject + mime + "\n" + r.body)
	err := smtp.SendMail(emailClient+":"+emailPort, auth, fromEmail, r.to, message)
	if err != nil {
		return false, err
	}
	return true, nil
}
