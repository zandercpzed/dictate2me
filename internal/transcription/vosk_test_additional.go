package transcription

import (
	"testing"
)

func TestParseVoskResult_Empty(t *testing.T) {
	seg, err := parseVoskResult(`{"text":""}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if seg.Text != "" {
		t.Fatalf("expected empty text, got %s", seg.Text)
	}
}
