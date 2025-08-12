package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	Sandbox = "sandbox"
	Prod    = "production"
)

// Recipient is a model
type Recipient struct {
	Number    string `json:"number"`
	Cost      string `json:"cost"`
	Status    string `json:"status"`
	MessageID string `json:"messageId"`
}

// SMS2 is a model
type SMS2 struct {
	Recipients []Recipient `json:"recipients"`
}

// SendMessageResponse is a model
type SendMessageResponse struct {
	SMS SMS2 `json:"SMSMessageData"`
}

// service is a model
type service struct {
	Username string
	APIKey   string
	URL      string
}

// NewSMS returns a new service
func NewSMS(username, apiKey, URL string) service {
	return service{username, apiKey, URL}
}

// Send - POST
func (s *service) Send(from, to, message string) (*SendMessageResponse, error) {
	values := url.Values{}
	values.Set("username", s.Username)
	values.Set("to", to)
	values.Set("message", message)
	if from != "" {
		// set from = "" to avoid this
		values.Set("from", from)
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	res, err := s.newPostRequest(s.URL, values, headers)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var smsMessageResponse SendMessageResponse
	if err := json.NewDecoder(res.Body).Decode(&smsMessageResponse); err != nil {
		return nil, errors.New("unable to parse sms response")
	}
	return &smsMessageResponse, nil
}

func (s *service) newPostRequest(url string, values url.Values, headers map[string]string) (*http.Response, error) {
	reader := strings.NewReader(values.Encode())

	req, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Content-Length", strconv.Itoa(reader.Len()))
	req.Header.Set("apikey", s.APIKey)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	return client.Do(req)
}
