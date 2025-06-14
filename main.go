package main

import (
	"github.com/mnako/letters"
	"io"
	"log"
	"mailhook/config"
	"time"

	"github.com/emersion/go-smtp"
)

func init() {
	if err := config.LoadConfig(); err != nil {
		log.Fatal(err)
	}
}

// The Backend implements SMTP server methods.
type Backend struct{}

// NewSession is called after client greeting (EHLO, HELO).
func (bkd *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &Session{}, nil
}

// A Session is returned after successful login.
type Session struct {
	from string
	to   string
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	// Check if the sender is allowed
	if !config.IsAllowedEmail(from) {
		log.Printf("Rejected mail from %s: not allowed", from)
		return smtp.ErrAuthFailed
	}

	s.from = from
	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	if !config.IsAllowedEmail(to) {
		log.Printf("Rejected mail to %s: not allowed", to)
		return smtp.ErrAuthFailed
	}
	s.to = to
	return nil
}

func (s *Session) Data(r io.Reader) error {
	log.Printf("Received email from %s to %s", s.from, s.to)

	subject := "No Subject"
	content := ""
	email, err := letters.ParseEmail(r)
	if err == nil {
		subject, content = email.Headers.Subject, email.Text
	} else {
		log.Printf("Error parsing email: %v", err)
	}

	if err := config.RunActions(s.from, s.to, subject, content); err != nil {
		log.Printf("Error running actions: %v", err)
	}
	return nil
}

func (s *Session) Reset() {
	s.from = ""
	s.to = ""
}

func (s *Session) Logout() error {
	return nil
}

func main() {

	be := &Backend{}

	s := smtp.NewServer(be)

	s.Addr = "0.0.0.0:2525"
	s.Domain = "localhost"
	s.WriteTimeout = 10 * time.Second
	s.ReadTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true

	log.Println("Starting server at", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
