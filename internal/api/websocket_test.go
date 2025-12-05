package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestWebsocket_AudioBeforeStart_ReturnsError(t *testing.T) {
	server, _ := setupTestServer(t)
	defer server.config.TranscriptionEngine.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/stream", server.handleStream)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http") + "/api/v1/stream"

	dialer := websocket.DefaultDialer
	conn, resp, err := dialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial failed: %v (resp: %v)", err, resp)
	}
	defer conn.Close()

	// Send audio message before start
	msg := StreamMessage{
		Type: "audio",
		Data: json.RawMessage(`{"data":"AA=="}`),
	}

	if err := conn.WriteJSON(msg); err != nil {
		t.Fatalf("write json failed: %v", err)
	}

	// Expect an error message back
	var got StreamMessage
	if err := conn.ReadJSON(&got); err != nil {
		t.Fatalf("read json failed: %v", err)
	}

	if got.Type != "error" {
		t.Fatalf("expected error message type, got %s", got.Type)
	}
}
