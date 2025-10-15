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

func TestRegister(t *testing.T) {
	client := CreateClient()
	registerId, config, err := client.Register(1, Hysteria2, "test-hostname", 8080, "127.0.0.1")
	if err != nil {
		t.Error(err)
	}
	t.Logf("RegisterId: %d, Config: %v", registerId, config)
}

func TestUsers(t *testing.T) {
	client := CreateClient()
	userList, err := client.Users(1, Hysteria2)
	t.Log(len(*userList))

	if err != nil && !errors.Is(err, ErrorUserNotModified) {
		t.Error(err)
	}
}

func TestUsers2(t *testing.T) {
	client := CreateClient()
	userList, err := client.Users(1, Hysteria2)
	t.Log(len(*userList))

	if err != nil {
		t.Error(err)
	}

	userList, err = client.Users(1, Hysteria2)
	t.Log(len(*userList))
	if err != nil && !errors.Is(err, ErrorUserNotModified) {
		t.Error(err)
	}
}

func TestSubmit(t *testing.T) {
	client := CreateClient()
	users, err := client.Users(1, Hysteria2)
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
	err = client.Submit(1, Hysteria2, generalUserTraffic)
	if err != nil {
		t.Error(err)
	}
}

func TestSubmitWithAgent(t *testing.T) {
	client := CreateClient()
	users, err := client.Users(1, Hysteria2)
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
	err = client.SubmitWithAgent(1, Hysteria2, generalUserTraffic)
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

	err := client.SubmitStatsWithAgent(1, Hysteria2, "127.0.0.1", stats)
	if err != nil {
		t.Error(err)
	}
}

func TestHeartbeat(t *testing.T) {
	client := CreateClient()
	err := client.Heartbeat(1, Hysteria2, "127.0.0.1")
	if err != nil {
		t.Error(err)
	}
}
