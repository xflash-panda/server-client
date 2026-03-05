package pkg

import (
	"context"
	"errors"
	"os"
	"testing"
)

func CreateClient() *Client {
	apiConfig := &Config{
		APIHost: "http://127.0.0.1:8080",
		Token:   "123456789123456789",
		Debug:   true,
	}
	client := New(apiConfig)
	return client
}

func skipIfNoServer(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("Skipping integration test; set INTEGRATION_TEST=1 to run")
	}
}

func TestIntegrationConfig(t *testing.T) {
	skipIfNoServer(t)
	client := CreateClient()
	ctx := context.Background()
	config, err := client.Config(ctx, 1, Trojan)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Config: %v", config)
}

func TestIntegrationRegister(t *testing.T) {
	skipIfNoServer(t)
	client := CreateClient()
	ctx := context.Background()
	registerId, err := client.Register(ctx, 32, Trojan, "test-hostname", 8080, "127.0.0.1")
	if err != nil {
		t.Error(err)
	}
	t.Logf("RegisterId: %s", registerId)
}

func TestIntegrationUsers(t *testing.T) {
	skipIfNoServer(t)
	client := CreateClient()
	ctx := context.Background()
	userList, err := client.Users(ctx, "1", Trojan)
	if err != nil && !errors.Is(err, ErrorUserNotModified) {
		t.Error(err)
	}
	if userList != nil {
		t.Log(len(*userList))
	}
}

func TestIntegrationUsers2(t *testing.T) {
	skipIfNoServer(t)
	client := CreateClient()
	ctx := context.Background()
	userList, err := client.Users(ctx, "1", Trojan)
	if err != nil {
		t.Error(err)
		return
	}
	if userList != nil {
		t.Log(len(*userList))
	}

	userList, err = client.Users(ctx, "1", Trojan)
	if err != nil && !errors.Is(err, ErrorUserNotModified) {
		t.Error(err)
	}
	if userList != nil {
		t.Log(len(*userList))
	}
}

func TestIntegrationSubmit(t *testing.T) {
	skipIfNoServer(t)
	client := CreateClient()
	ctx := context.Background()
	users, err := client.Users(ctx, "1", Trojan)
	if err != nil {
		t.Error(err)
		return
	}
	generalUserTraffic := make([]*UserTraffic, len(*users))
	for i, userInfo := range *users {
		generalUserTraffic[i] = &UserTraffic{
			UID:      userInfo.ID,
			Upload:   114514,
			Download: 114514,
			Count:    33,
		}
	}
	err = client.Submit(ctx, "1", Trojan, generalUserTraffic)
	if err != nil {
		t.Error(err)
	}
}

func TestIntegrationSubmitWithAgent(t *testing.T) {
	skipIfNoServer(t)
	client := CreateClient()
	ctx := context.Background()
	users, err := client.Users(ctx, "1", Trojan)
	if err != nil {
		t.Error(err)
		return
	}
	generalUserTraffic := make([]*UserTraffic, len(*users))
	for i, userInfo := range *users {
		generalUserTraffic[i] = &UserTraffic{
			UID:      userInfo.ID,
			Upload:   114514,
			Download: 114514,
			Count:    22,
		}
	}
	err = client.SubmitWithAgent(ctx, "1", Trojan, generalUserTraffic)
	if err != nil {
		t.Error(err)
	}
}

func TestIntegrationSubmitStatsWithAgent(t *testing.T) {
	skipIfNoServer(t)
	client := CreateClient()
	ctx := context.Background()
	stats := &TrafficStats{
		Count:    1,
		Requests: 1,
		UserIds:  []int{1, 2, 3},
		UserRequests: map[int]int{
			1: 2, 3: 4,
		},
	}

	err := client.SubmitStatsWithAgent(ctx, "1", Trojan, stats)
	if err != nil {
		t.Error(err)
	}
}

func TestIntegrationHeartbeat(t *testing.T) {
	skipIfNoServer(t)
	client := CreateClient()
	ctx := context.Background()
	err := client.Heartbeat(ctx, "1", Trojan, "127.0.0.1")
	if err != nil {
		t.Error(err)
	}
}

func TestIntegrationVerify(t *testing.T) {
	skipIfNoServer(t)
	client := CreateClient()
	ctx := context.Background()
	valid, err := client.Verify(ctx, "1", Trojan)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Verify result: %v", valid)
}
