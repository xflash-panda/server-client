package pkg

import (
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

func TestConfig(t *testing.T) {
	skipIfNoServer(t)
	client := CreateClient()
	config, err := client.Config(1, Trojan)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Config: %v", config)
}

func TestRegister(t *testing.T) {
	skipIfNoServer(t)
	client := CreateClient()
	registerId, err := client.Register(32, Trojan, "test-hostname", 8080, "127.0.0.1")
	if err != nil {
		t.Error(err)
	}
	t.Logf("RegisterId: %s", registerId)
}

func TestUsers(t *testing.T) {
	skipIfNoServer(t)
	client := CreateClient()
	userList, err := client.Users("1", Trojan)
	if err != nil && !errors.Is(err, ErrorUserNotModified) {
		t.Error(err)
	}
	if userList != nil {
		t.Log(len(*userList))
	}
}

func TestUsers2(t *testing.T) {
	skipIfNoServer(t)
	client := CreateClient()
	userList, err := client.Users("1", Trojan)
	if err != nil {
		t.Error(err)
		return
	}
	if userList != nil {
		t.Log(len(*userList))
	}

	userList, err = client.Users("1", Trojan)
	if err != nil && !errors.Is(err, ErrorUserNotModified) {
		t.Error(err)
	}
	if userList != nil {
		t.Log(len(*userList))
	}
}

func TestSubmit(t *testing.T) {
	skipIfNoServer(t)
	client := CreateClient()
	users, err := client.Users("1", Trojan)
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
	err = client.Submit("1", Trojan, generalUserTraffic)
	if err != nil {
		t.Error(err)
	}
}

func TestSubmitWithAgent(t *testing.T) {
	skipIfNoServer(t)
	client := CreateClient()
	users, err := client.Users("1", Trojan)
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
	err = client.SubmitWithAgent("1", Trojan, generalUserTraffic)
	if err != nil {
		t.Error(err)
	}
}

func TestSubmitStatsWithAgent(t *testing.T) {
	skipIfNoServer(t)
	client := CreateClient()
	stats := &TrafficStats{
		Count:    1,
		Requests: 1,
		UserIds:  []int{1, 2, 3},
		UserRequests: map[int]int{
			1: 2, 3: 4,
		},
	}

	err := client.SubmitStatsWithAgent("1", Trojan, stats)
	if err != nil {
		t.Error(err)
	}
}

func TestHeartbeat(t *testing.T) {
	skipIfNoServer(t)
	client := CreateClient()
	err := client.Heartbeat("1", Trojan, "127.0.0.1")
	if err != nil {
		t.Error(err)
	}
}

func TestVerify(t *testing.T) {
	skipIfNoServer(t)
	client := CreateClient()
	valid, err := client.Verify("1", Trojan)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Verify result: %v", valid)
}
