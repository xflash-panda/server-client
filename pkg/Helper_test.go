package pkg

import (
	"testing"
)

func TestAsHysteria2(t *testing.T) {
	client := CreateClient()
	client.Debug(false)
	config, err := client.Config(1, Hysteria2)
	if err != nil {
		t.Error(err)
	}

	//hy2Config := config.(*Hysteria2Config)
	hy2Config, err := AsHysteria2Config(config)
	if err != nil {
		t.Error(err)
	}

	t.Log(hy2Config.ID)
}

func TestUnmarshalHysteriaConfig(t *testing.T) {
	client := CreateClient()
	client.Debug(false)
	configBytes, err := client.RawConfig(1, Hysteria2)
	if err != nil {
		t.Error(err)
	}

	hy2Config, err := UnmarshalHysteria2Config(configBytes)
	if err != nil {
		t.Error(err)
	}
	t.Log(hy2Config.ID)
}

func TestUnmarshalUsers(t *testing.T) {
	client := CreateClient()
	client.Debug(false)
	usersBytes, err := client.RawUsers(1, Hysteria2)
	if err != nil {
		t.Error(err)
	}

	users, err := UnmarshalUsers(usersBytes)
	if err != nil {
		t.Error(err)
	}
	t.Log(len(*users))
}
