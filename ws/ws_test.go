package ws

import (
	"strings"
	"testing"
	"time"
)

func TestToWSURL(t *testing.T) {
	// Test HTTP -> WS
	url1, err := toWSURL("http://localhost:8080", "console", "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(url1, "ws://localhost:8080/v1/ws/console?key=test-key") {
		t.Errorf("expected ws:// url with key, got %s", url1)
	}

	// Test HTTPS -> WSS
	url2, err := toWSURL("https://api.example.com/", "/events/", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if url2 != "wss://api.example.com/v1/ws/events" {
		t.Errorf("expected wss://api.example.com/v1/ws/events, got %s", url2)
	}
}

func TestBackoffDelay(t *testing.T) {
	max := 10 * time.Second
	for i := 0; i < 5; i++ {
		d := backoffDelay(i, max)
		if d < 0 || d > max {
			t.Errorf("backoff delay out of range: %v", d)
		}
	}
}
