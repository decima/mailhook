package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Email struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
}

type Handler struct {
	From   []string `json:"from"`
	To     []string `json:"to"`
	Action Action   `json:"action"`
}

type Action struct {
	URL    string         `json:"url"`
	Method string         `json:"method"`
	Format string         `json:"format"`
	Body   map[string]any `json:"body,omitempty"`
}

func (a *Action) Run(e Email) error {
	httpClient := &http.Client{}

	req, err := http.NewRequest(a.Method, a.URL, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	var content string = ""
	if a.Format != "empty" {
		for k, v := range map[string]string{
			"from":    e.From,
			"to":      e.To,
			"subject": e.Subject,
			"text":    e.Text,
		} {
			a.Body["_"+k] = v
		}
		switch a.Format {
		case "json":
			jsonBody, err := json.Marshal(a.Body)
			if err != nil {
				return fmt.Errorf("error marshalling body to JSON: %w", err)
			}
			content = string(jsonBody)
			req.Header.Set("Content-Type", "application/json")
		default:
			return fmt.Errorf("unsupported format: %s", a.Format)
		}

		req.Body = io.NopCloser(bytes.NewBufferString(content))
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}

	log.Println("Response status:", resp.Status)
	defer resp.Body.Close()

	return nil
}

type Allowed struct {
	Emails []string
}

func (allowed *Allowed) Add(email string) {
	if email != "" && !allowed.Has(email) {
		allowed.Emails = append(allowed.Emails, email)
	}
}

func (allowed *Allowed) Has(email string) bool {
	for _, e := range allowed.Emails {
		if e == email {
			return true
		}
	}
	return false
}

func (allowed *Allowed) IsAllowed(email string) bool {
	for _, e := range allowed.Emails {
		if e == email {
			return true
		}
	}
	return false
}

var Recipients = Allowed{}
var Senders = Allowed{}

func LoadConfig() error {
	var cfg []Handler
	fileName := "config.json"
	data, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("error reading config file %s: %w", fileName, err)
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("error parsing config file %s: %w", fileName, err)
	}
	log.Printf("Loaded config")

	for _, h := range cfg {
		for _, a := range h.From {
			Senders.Add(a)
		}
		for _, a := range h.To {
			Recipients.Add(a)
		}
		for _, from := range h.From {
			for _, to := range h.To {
				fromToPair := uniquePair(from, to)
				Config[fromToPair] = append(Config[fromToPair], h.Action)
			}
		}
	}

	return nil
}

var Config = map[string][]Action{}

func uniquePair(from string, to string) string {
	return fmt.Sprintf("%s->%s", from, to)
}

func IsAllowedEmail(email string) bool {
	return Recipients.IsAllowed(email) || Senders.IsAllowed(email)
}
func RunActions(from string, to string, subject string, text string) error {
	email := Email{
		From:    from,
		To:      to,
		Subject: subject,
		Text:    text,
	}
	fromToPair := uniquePair(from, to)
	actions, exists := Config[fromToPair]
	if !exists {
		log.Printf("No actions for %s from %s to %s", subject, from, to)
		return nil
	}
	for _, action := range actions {
		if err := action.Run(email); err != nil {
			return err
		}
	}
	return nil
}
