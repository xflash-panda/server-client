package pkg

import (
	"errors"
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

func TestConfig(t *testing.T) {
	client := CreateClient()
	config, err := client.Config(1, Trojan)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Config: %v", config)
}

func TestRegister(t *testing.T) {
	client := CreateClient()
	registerId, err := client.Register(32, Trojan, "test-hostname", 8080, "127.0.0.1")
	if err != nil {
		t.Error(err)
	}
	t.Logf("RegisterId: %s", registerId)
}

func TestUsers(t *testing.T) {
	client := CreateClient()
	userList, err := client.Users("1", Trojan)
	t.Log(len(*userList))

	if err != nil && !errors.Is(err, ErrorUserNotModified) {
		t.Error(err)
	}
}

func TestUsers2(t *testing.T) {
	client := CreateClient()
	userList, err := client.Users("1", Trojan)
	t.Log(len(*userList))

	if err != nil {
		t.Error(err)
	}

	userList, err = client.Users("1", Trojan)
	t.Log(len(*userList))
	if err != nil && !errors.Is(err, ErrorUserNotModified) {
		t.Error(err)
	}
}

func TestSubmit(t *testing.T) {
	client := CreateClient()
	users, err := client.Users("1", Trojan)
	if err != nil {
		t.Error(err)
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
	// client.Debug()
	err = client.Submit("1", Trojan, generalUserTraffic)
	if err != nil {
		t.Error(err)
	}
}

func TestSubmitWithAgent(t *testing.T) {
	client := CreateClient()
	users, err := client.Users("1", Trojan)
	if err != nil {
		t.Error(err)
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
	client := CreateClient()
	stats := &TrafficStats{
		Count:    1,
		Requests: 1,
		UserIds:  []int{1, 2, 3},
		UserRequests: map[int]int{
			1: 2, 3: 4,
		},
	}

	err := client.SubmitStatsWithAgent("1", Trojan, "127.0.0.1", stats)
	if err != nil {
		t.Error(err)
	}
}

func TestHeartbeat(t *testing.T) {
	client := CreateClient()
	err := client.Heartbeat("1", Trojan, "127.0.0.1")
	if err != nil {
		t.Error(err)
	}
}
