package integration

import (
	"bytes"
	"net/http"
	"testing"
)

func TestEventEndpoint(t *testing.T) {
	eventJSON := `{"type":"UserCreated","payload":{"id":1}}`

	resp, err := http.Post("http://localhost:8080/events", "application/json", bytes.NewBuffer([]byte(eventJSON)))
	if err != nil {
		t.Fatalf("Failed to POST event: %v", err)
	}

	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("Expected 202 Accepted, got %d", resp.StatusCode)
	}
}
