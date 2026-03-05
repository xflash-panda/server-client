package pkg

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// newTestServer creates a test HTTP server that returns the given status code and body.
func newTestServer(t *testing.T, statusCode int, body any) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		if body != nil {
			_ = json.NewEncoder(w).Encode(body)
		}
	}))
	t.Cleanup(server.Close)
	return server
}

// newTestClient creates a Client pointing to the test server.
func newTestClient(t *testing.T, serverURL string) *Client {
	t.Helper()
	return New(&Config{
		APIHost: serverURL,
		Token:   "test-token",
		Timeout: 5 * time.Second,
	})
}

func TestConfig(t *testing.T) {
	resp := map[string]any{
		"data": map[string]any{
			"id":          1,
			"server_port": 443,
			"server_name": "example.com",
			"network":     "tcp",
		},
		"message": "success",
	}
	server := newTestServer(t, 200, resp)
	client := newTestClient(t, server.URL)

	ctx := context.Background()
	config, err := client.Config(ctx, 1, Trojan)
	if err != nil {
		t.Fatalf("Config() unexpected error: %v", err)
	}
	if config == nil {
		t.Fatal("Config() returned nil config")
	}
}

func TestRawConfig(t *testing.T) {
	resp := map[string]any{
		"data": map[string]any{
			"id":          1,
			"server_port": 443,
		},
		"message": "success",
	}
	server := newTestServer(t, 200, resp)
	client := newTestClient(t, server.URL)

	ctx := context.Background()
	rawData, err := client.RawConfig(ctx, 1, Trojan)
	if err != nil {
		t.Fatalf("RawConfig() unexpected error: %v", err)
	}
	if len(rawData) == 0 {
		t.Fatal("RawConfig() returned empty data")
	}
}

func TestRegister(t *testing.T) {
	resp := map[string]any{
		"data": map[string]any{
			"register_id": "test-register-id",
		},
		"message": "success",
	}
	server := newTestServer(t, 200, resp)
	client := newTestClient(t, server.URL)

	ctx := context.Background()
	registerId, err := client.Register(ctx, 1, Trojan, "test-hostname", 8080, "127.0.0.1")
	if err != nil {
		t.Fatalf("Register() unexpected error: %v", err)
	}
	if registerId != "test-register-id" {
		t.Fatalf("Expected register_id='test-register-id', got '%s'", registerId)
	}
}

func TestUnregister(t *testing.T) {
	resp := map[string]any{
		"data":    true,
		"message": "success",
	}
	server := newTestServer(t, 200, resp)
	client := newTestClient(t, server.URL)

	ctx := context.Background()
	err := client.Unregister(ctx, Trojan, "test-register-id")
	if err != nil {
		t.Fatalf("Unregister() unexpected error: %v", err)
	}
}

func TestRawUsers(t *testing.T) {
	resp := map[string]any{
		"data": []map[string]any{
			{"id": 1, "uuid": "uuid-1"},
			{"id": 2, "uuid": "uuid-2"},
		},
		"message": "success",
	}
	server := newTestServer(t, 200, resp)
	client := newTestClient(t, server.URL)

	ctx := context.Background()
	rawData, err := client.RawUsers(ctx, "test-register-id", Trojan)
	if err != nil {
		t.Fatalf("RawUsers() unexpected error: %v", err)
	}
	if len(rawData) == 0 {
		t.Fatal("RawUsers() returned empty data")
	}
}

func TestUsers(t *testing.T) {
	resp := map[string]any{
		"data": []map[string]any{
			{"id": 1, "uuid": "uuid-1"},
			{"id": 2, "uuid": "uuid-2"},
		},
		"message": "success",
	}
	server := newTestServer(t, 200, resp)
	client := newTestClient(t, server.URL)

	ctx := context.Background()
	users, err := client.Users(ctx, "test-register-id", Trojan)
	if err != nil {
		t.Fatalf("Users() unexpected error: %v", err)
	}
	if users == nil || len(*users) != 2 {
		t.Fatalf("Expected 2 users, got %v", users)
	}
}

func TestRawUsersByNodeId(t *testing.T) {
	resp := map[string]any{
		"data": []map[string]any{
			{"id": 1, "uuid": "uuid-1"},
		},
		"message": "success",
	}
	server := newTestServer(t, 200, resp)
	client := newTestClient(t, server.URL)

	ctx := context.Background()
	rawData, err := client.RawUsersByNodeId(ctx, 1, Trojan)
	if err != nil {
		t.Fatalf("RawUsersByNodeId() unexpected error: %v", err)
	}
	if len(rawData) == 0 {
		t.Fatal("RawUsersByNodeId() returned empty data")
	}
}

func TestUsersByNodeId(t *testing.T) {
	resp := map[string]any{
		"data": []map[string]any{
			{"id": 1, "uuid": "uuid-1"},
		},
		"message": "success",
	}
	server := newTestServer(t, 200, resp)
	client := newTestClient(t, server.URL)

	ctx := context.Background()
	users, err := client.UsersByNodeId(ctx, 1, Trojan)
	if err != nil {
		t.Fatalf("UsersByNodeId() unexpected error: %v", err)
	}
	if users == nil || len(*users) != 1 {
		t.Fatalf("Expected 1 user, got %v", users)
	}
}

func TestSubmit(t *testing.T) {
	resp := map[string]any{
		"data":    true,
		"message": "success",
	}
	server := newTestServer(t, 200, resp)
	client := newTestClient(t, server.URL)

	ctx := context.Background()
	traffic := []*UserTraffic{
		{UID: 1, Upload: 1024, Download: 2048, Count: 1},
	}
	err := client.Submit(ctx, "test-register-id", Trojan, traffic)
	if err != nil {
		t.Fatalf("Submit() unexpected error: %v", err)
	}
}

func TestSubmitWithAgent(t *testing.T) {
	resp := map[string]any{
		"data":    true,
		"message": "success",
	}
	server := newTestServer(t, 200, resp)
	client := newTestClient(t, server.URL)

	ctx := context.Background()
	traffic := []*UserTraffic{
		{UID: 1, Upload: 1024, Download: 2048, Count: 1},
	}
	err := client.SubmitWithAgent(ctx, "test-register-id", Trojan, traffic)
	if err != nil {
		t.Fatalf("SubmitWithAgent() unexpected error: %v", err)
	}
}

func TestSubmitStatsWithAgent(t *testing.T) {
	resp := map[string]any{
		"data":    true,
		"message": "success",
	}
	server := newTestServer(t, 200, resp)
	client := newTestClient(t, server.URL)

	ctx := context.Background()
	stats := &TrafficStats{
		Count:    1,
		Requests: 1,
		UserIds:  []int{1},
		UserRequests: map[int]int{
			1: 1,
		},
	}
	err := client.SubmitStatsWithAgent(ctx, "test-register-id", Trojan, stats)
	if err != nil {
		t.Fatalf("SubmitStatsWithAgent() unexpected error: %v", err)
	}
}

func TestHeartbeat(t *testing.T) {
	resp := map[string]any{
		"data":    true,
		"message": "success",
	}
	server := newTestServer(t, 200, resp)
	client := newTestClient(t, server.URL)

	ctx := context.Background()
	err := client.Heartbeat(ctx, "test-register-id", Trojan, "127.0.0.1")
	if err != nil {
		t.Fatalf("Heartbeat() unexpected error: %v", err)
	}
}

func TestVerify(t *testing.T) {
	resp := map[string]any{
		"data":    true,
		"message": "success",
	}
	server := newTestServer(t, 200, resp)
	client := newTestClient(t, server.URL)

	ctx := context.Background()
	valid, err := client.Verify(ctx, "test-register-id", Trojan)
	if err != nil {
		t.Fatalf("Verify() unexpected error: %v", err)
	}
	if !valid {
		t.Fatal("Expected valid=true")
	}
}

func TestContextCancellation(t *testing.T) {
	// Server that delays response to allow cancellation
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			return
		case <-time.After(5 * time.Second):
			w.WriteHeader(200)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data":    true,
				"message": "success",
			})
		}
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, server.URL)

	ctx, cancel := context.WithCancel(context.Background())
	// Cancel immediately
	cancel()

	_, err := client.RawConfig(ctx, 1, Trojan)
	if err == nil {
		t.Fatal("Expected error from cancelled context, got nil")
	}
}

func TestContextTimeout(t *testing.T) {
	// Server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			return
		case <-time.After(5 * time.Second):
			w.WriteHeader(200)
		}
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, server.URL)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := client.RawConfig(ctx, 1, Trojan)
	if err == nil {
		t.Fatal("Expected error from context timeout, got nil")
	}
}

func TestHTTPError(t *testing.T) {
	server := newTestServer(t, 500, map[string]any{
		"message": "internal server error",
	})
	client := newTestClient(t, server.URL)

	ctx := context.Background()
	_, err := client.Config(ctx, 1, Trojan)
	if err == nil {
		t.Fatal("Expected error from 500 status code, got nil")
	}
}
