package notifier

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebhookNotifier(t *testing.T) {
	var receivedPayload map[string]string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}
		var payload map[string]string
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		receivedPayload = payload
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := NewWebhookNotifier(server.URL)

	err := n.Send(context.Background(), "Test Subject", "Test Body")
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}

	if receivedPayload["subject"] != "Test Subject" {
		t.Errorf("subject = %q, want %q", receivedPayload["subject"], "Test Subject")
	}
	if receivedPayload["body"] != "Test Body" {
		t.Errorf("body = %q, want %q", receivedPayload["body"], "Test Body")
	}
	if receivedPayload["source"] != "NetConfigHub" {
		t.Errorf("source = %q, want %q", receivedPayload["source"], "NetConfigHub")
	}
	if _, ok := receivedPayload["timestamp"]; !ok {
		t.Error("timestamp field missing in payload")
	}
}

func TestWebhookNotifierServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n := NewWebhookNotifier(server.URL)

	err := n.Send(context.Background(), "Subject", "Body")
	if err == nil {
		t.Error("expected error for non-2xx status, got nil")
	}
}

func TestWebhookNotifierInvalidURL(t *testing.T) {
	n := NewWebhookNotifier("http://invalid-host-that-does-not-exist.local:99999")

	err := n.Send(context.Background(), "Subject", "Body")
	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestNewNotifier(t *testing.T) {
	tests := []struct {
		name    string
		config  NotifyConfig
		wantNil bool
	}{
		{
			name:    "disabled",
			config:  NotifyConfig{Enabled: false},
			wantNil: true,
		},
		{
			name:    "empty type",
			config:  NotifyConfig{Enabled: true, Type: ""},
			wantNil: true,
		},
		{
			name:    "webhook type",
			config:  NotifyConfig{Enabled: true, Type: "webhook", WebhookURL: "http://example.com"},
			wantNil: false,
		},
		{
			name:    "email type",
			config:  NotifyConfig{Enabled: true, Type: "email", SMTPHost: "smtp.example.com"},
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NewNotifier(tt.config)
			if (n == nil) != tt.wantNil {
				t.Errorf("NewNotifier() returned nil = %v, want nil = %v", n == nil, tt.wantNil)
			}
		})
	}
}
