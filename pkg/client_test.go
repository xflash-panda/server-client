package pkg

import (
	"testing"
)

func CreateClient() *Client {
	apiConfig := &Config{
		APIHost:  "http://127.0.0.1:8080",
		NodeID:   1,
		Token:    "123456789123456789",
		NodeType: Hysteria2,
		Debug:    true,
	}
	client := New(apiConfig)
	return client
}

func TestConfig(t *testing.T) {
	client := CreateClient()
	config, err := client.Config()
	if err != nil {
		t.Error(err)
	}
	t.Log(config)
}

func TestUsers(t *testing.T) {
	client := CreateClient()
	userList, err := client.Users()
	if err != nil {
		t.Error(err)
	}
	t.Log(userList)
}

func TestSubmit(t *testing.T) {
	client := CreateClient()
	users, err := client.Users()
	if err != nil {
		t.Error(err)
	}
	generalUserTraffic := make([]*UserTraffic, len(*users))
	for i, userInfo := range *users {
		generalUserTraffic[i] = &UserTraffic{
			UID:      userInfo.ID,
			Upload:   114514,
			Download: 114514,
		}
	}
	//client.Debug()
	err = client.Submit(generalUserTraffic)
	if err != nil {
		t.Error(err)
	}
}

func TestSubmitWithAgent(t *testing.T) {
	client := CreateClient()
	users, err := client.Users()
	if err != nil {
		t.Error(err)
	}
	generalUserTraffic := make([]*UserTraffic, len(*users))
	for i, userInfo := range *users {
		generalUserTraffic[i] = &UserTraffic{
			UID:      userInfo.ID,
			Upload:   114514,
			Download: 114514,
		}
	}
	err = client.SubmitWithAgent("127.0.0.1", generalUserTraffic)
	if err != nil {
		t.Error(err)
	}
}

func TestSubmitStatsWithAgent(t *testing.T) {
	client := CreateClient()
	stats := &Stats{
		Count:    1,
		Requests: 1,
		UserIds:  []int{1, 2, 3},
		UserRequests: map[int]int{
			1: 2, 3: 4,
		},
	}

	err := client.SubmitStatsWithAgent("127.0.0.1", stats)
	if err != nil {
		t.Error(err)
	}
}

func TestHeartbeat(t *testing.T) {
	client := CreateClient()
	err := client.Heartbeat("127.0.0.1")
	if err != nil {
		t.Error(err)
	}
}
