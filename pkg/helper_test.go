package pkg

import (
	"encoding/json"
	"testing"

	"github.com/xflash-panda/server-client/pkg/xray"
)

// TestAsConfig 测试泛型类型转换函数
func TestAsConfig(t *testing.T) {
	tests := []struct {
		name    string
		input   NodeConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid VMess config conversion",
			input: &VMessConfig{
				ID:         1,
				ServerPort: 8080,
				TLS:        1,
				Network:    "tcp",
			},
			wantErr: false,
		},
		{
			name: "valid Hysteria2 config conversion",
			input: &Hysteria2Config{
				ID:         2,
				ServerPort: 9090,
				Obfs:       "salamander",
				UpMbps:     100,
				DownMbps:   200,
			},
			wantErr: false,
		},
		{
			name:    "nil pointer conversion",
			input:   (*VMessConfig)(nil),
			wantErr: true,
			errMsg:  "nil cannot be converted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch cfg := tt.input.(type) {
			case *VMessConfig:
				result, err := AsConfig[*VMessConfig](cfg)
				if tt.wantErr {
					if err == nil {
						t.Errorf("AsConfig() expected error but got nil")
					}
				} else {
					if err != nil {
						t.Errorf("AsConfig() unexpected error: %v", err)
					}
					if result == nil {
						t.Errorf("AsConfig() expected non-nil result")
					}
				}
			case *Hysteria2Config:
				result, err := AsConfig[*Hysteria2Config](cfg)
				if tt.wantErr {
					if err == nil {
						t.Errorf("AsConfig() expected error but got nil")
					}
				} else {
					if err != nil {
						t.Errorf("AsConfig() unexpected error: %v", err)
					}
					if result == nil {
						t.Errorf("AsConfig() expected non-nil result")
					}
				}
			}
		})
	}
}

// TestAsConfigTypeAssertionFailure 测试类型断言失败的情况
func TestAsConfigTypeAssertionFailure(t *testing.T) {
	// 尝试将 VMess 配置转换为 Hysteria2 配置（应该失败）
	vMessConfig := &VMessConfig{
		ID:         1,
		ServerPort: 8080,
	}

	result, err := AsConfig[*Hysteria2Config](vMessConfig)
	if err == nil {
		t.Errorf("AsConfig() expected error when converting VMess to Hysteria2, but got nil")
	}
	if result != nil {
		t.Errorf("AsConfig() expected nil result on error, but got %v", result)
	}
}

// TestAsVMessConfig 测试 VMess 配置转换
func TestAsVMessConfig(t *testing.T) {
	tests := []struct {
		name    string
		input   NodeConfig
		wantErr bool
	}{
		{
			name: "valid VMess config",
			input: &VMessConfig{
				ID:         1,
				ServerPort: 8080,
				TLS:        1,
				Network:    "tcp",
			},
			wantErr: false,
		},
		{
			name: "invalid config type",
			input: &Hysteria2Config{
				ID: 2,
			},
			wantErr: true,
		},
		{
			name:    "nil config",
			input:   (*VMessConfig)(nil),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := AsVMessConfig(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("AsVMessConfig() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("AsVMessConfig() unexpected error: %v", err)
				}
				if result == nil {
					t.Errorf("AsVMessConfig() expected non-nil result")
				}
			}
		})
	}
}

// TestAsHysteriaConfig 测试 Hysteria 配置转换
func TestAsHysteriaConfig(t *testing.T) {
	tests := []struct {
		name    string
		input   NodeConfig
		wantErr bool
	}{
		{
			name: "valid Hysteria config",
			input: &HysteriaConfig{
				ID:         1,
				ServerPort: 8080,
				Protocol:   "udp",
				Obfs:       "xplus",
				UpMbps:     100,
				DownMbps:   200,
			},
			wantErr: false,
		},
		{
			name: "invalid config type",
			input: &VMessConfig{
				ID: 1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := AsHysteriaConfig(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("AsHysteriaConfig() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("AsHysteriaConfig() unexpected error: %v", err)
				}
				if result == nil {
					t.Errorf("AsHysteriaConfig() expected non-nil result")
				}
			}
		})
	}
}

// TestAsHysteria2Config 测试 Hysteria2 配置转换
func TestAsHysteria2Config(t *testing.T) {
	tests := []struct {
		name    string
		input   NodeConfig
		wantErr bool
	}{
		{
			name: "valid Hysteria2 config",
			input: &Hysteria2Config{
				ID:                 1,
				ServerPort:         8080,
				Obfs:               "salamander",
				UpMbps:             100,
				DownMbps:           200,
				IgnoreCliBandWidth: true,
				DisableUDP:         false,
			},
			wantErr: false,
		},
		{
			name: "invalid config type",
			input: &TrojanConfig{
				ID: 1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := AsHysteria2Config(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("AsHysteria2Config() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("AsHysteria2Config() unexpected error: %v", err)
				}
				if result == nil {
					t.Errorf("AsHysteria2Config() expected non-nil result")
				}
			}
		})
	}
}

// TestAsTrojanConfig 测试 Trojan 配置转换
func TestAsTrojanConfig(t *testing.T) {
	tests := []struct {
		name    string
		input   NodeConfig
		wantErr bool
	}{
		{
			name: "valid Trojan config",
			input: &TrojanConfig{
				ID:            1,
				ServerPort:    443,
				AllowInsecure: 0,
				ServerName:    "example.com",
				Network:       "tcp",
			},
			wantErr: false,
		},
		{
			name: "valid Trojan config with websocket",
			input: &TrojanConfig{
				ID:         1,
				ServerPort: 443,
				Network:    "ws",
				WebSocketConfig: &xray.WebSocketConfig{
					Path: "/ws",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid config type",
			input: &ShadowsocksConfig{
				ID: 1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := AsTrojanConfig(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("AsTrojanConfig() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("AsTrojanConfig() unexpected error: %v", err)
				}
				if result == nil {
					t.Errorf("AsTrojanConfig() expected non-nil result")
				}
			}
		})
	}
}

// TestAsShadowsocksConfig 测试 Shadowsocks 配置转换
func TestAsShadowsocksConfig(t *testing.T) {
	tests := []struct {
		name    string
		input   NodeConfig
		wantErr bool
	}{
		{
			name: "valid Shadowsocks config",
			input: &ShadowsocksConfig{
				ID:         1,
				ServerPort: 8388,
				Method:     "aes-256-gcm",
				Network:    "tcp,udp",
			},
			wantErr: false,
		},
		{
			name: "invalid config type",
			input: &AnyTLSConfig{
				ID: 1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := AsShadowsocksConfig(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("AsShadowsocksConfig() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("AsShadowsocksConfig() unexpected error: %v", err)
				}
				if result == nil {
					t.Errorf("AsShadowsocksConfig() expected non-nil result")
				}
			}
		})
	}
}

// TestAsAnyTLSConfig 测试 AnyTLS 配置转换
func TestAsAnyTLSConfig(t *testing.T) {
	tests := []struct {
		name    string
		input   NodeConfig
		wantErr bool
	}{
		{
			name: "valid AnyTLS config",
			input: &AnyTLSConfig{
				ID:            1,
				ServerPort:    443,
				AllowInsecure: 0,
				ServerName:    "example.com",
				PaddingRules:  "random",
			},
			wantErr: false,
		},
		{
			name: "invalid config type",
			input: &VMessConfig{
				ID: 1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := AsAnyTLSConfig(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("AsAnyTLSConfig() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("AsAnyTLSConfig() unexpected error: %v", err)
				}
				if result == nil {
					t.Errorf("AsAnyTLSConfig() expected non-nil result")
				}
			}
		})
	}
}

// TestUnmarshalConfig 测试泛型反序列化函数
func TestUnmarshalConfig(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		wantErr bool
	}{
		{
			name: "valid Hysteria2 JSON",
			input: []byte(`{
				"data": {
					"id": 1,
					"server_port": 8080,
					"obfs": "salamander",
					"up_mbps": 100,
					"down_mbps": 200,
					"ignore_cli_band_width": true,
					"disable_udp": false
				},
				"message": "success"
			}`),
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			input:   []byte(`{invalid json}`),
			wantErr: true,
		},
		{
			name:    "empty data",
			input:   []byte(`{"data": null, "message": "success"}`),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := UnmarshalConfig[Hysteria2Config](tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("UnmarshalConfig() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("UnmarshalConfig() unexpected error: %v", err)
				}
			}
			// 注意：即使成功解析，result 也可能为 nil（当 data 字段为 null 时）
			_ = result
		})
	}
}

// TestUnmarshalHysteria2Config 测试 Hysteria2 反序列化
func TestUnmarshalHysteria2Config(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		wantErr  bool
		validate func(*testing.T, *Hysteria2Config)
	}{
		{
			name: "valid Hysteria2 config",
			input: []byte(`{
				"data": {
					"id": 1,
					"server_port": 8080,
					"obfs": "salamander",
					"up_mbps": 100,
					"down_mbps": 200,
					"ignore_cli_band_width": true,
					"disable_udp": false
				},
				"message": "success"
			}`),
			wantErr: false,
			validate: func(t *testing.T, config *Hysteria2Config) {
				if config == nil {
					t.Errorf("Expected non-nil config")
					return
				}
				if config.ID != 1 {
					t.Errorf("Expected ID=1, got %d", config.ID)
				}
				if config.ServerPort != 8080 {
					t.Errorf("Expected ServerPort=8080, got %d", config.ServerPort)
				}
				if config.Obfs != "salamander" {
					t.Errorf("Expected Obfs=salamander, got %s", config.Obfs)
				}
				if config.UpMbps != 100 {
					t.Errorf("Expected UpMbps=100, got %d", config.UpMbps)
				}
				if config.DownMbps != 200 {
					t.Errorf("Expected DownMbps=200, got %d", config.DownMbps)
				}
				if !config.IgnoreCliBandWidth {
					t.Errorf("Expected IgnoreCliBandWidth=true")
				}
				if config.DisableUDP {
					t.Errorf("Expected DisableUDP=false")
				}
			},
		},
		{
			name:    "invalid JSON",
			input:   []byte(`{invalid}`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := UnmarshalHysteria2Config(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("UnmarshalHysteria2Config() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("UnmarshalHysteria2Config() unexpected error: %v", err)
				}
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}
		})
	}
}

// TestUnmarshalHysteriaConfig 测试 Hysteria 反序列化
func TestUnmarshalHysteriaConfig(t *testing.T) {
	input := []byte(`{
		"data": {
			"id": 1,
			"server_port": 8080,
			"protocol": "udp",
			"obfs": "xplus",
			"up_mbps": 100,
			"down_mbps": 200,
			"disable_mtu_discovery": true,
			"disable_udp": false
		},
		"message": "success"
	}`)

	result, err := UnmarshalHysteriaConfig(input)
	if err != nil {
		t.Fatalf("UnmarshalHysteriaConfig() unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if result.ID != 1 {
		t.Errorf("Expected ID=1, got %d", result.ID)
	}
	if result.ServerPort != 8080 {
		t.Errorf("Expected ServerPort=8080, got %d", result.ServerPort)
	}
	if result.Protocol != "udp" {
		t.Errorf("Expected Protocol=udp, got %s", result.Protocol)
	}
}

// TestUnmarshalTrojanConfig 测试 Trojan 反序列化
func TestUnmarshalTrojanConfig(t *testing.T) {
	input := []byte(`{
		"data": {
			"id": 1,
			"server_port": 443,
			"allow_insecure": 0,
			"server_name": "example.com",
			"network": "tcp"
		},
		"message": "success"
	}`)

	result, err := UnmarshalTrojanConfig(input)
	if err != nil {
		t.Fatalf("UnmarshalTrojanConfig() unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if result.ID != 1 {
		t.Errorf("Expected ID=1, got %d", result.ID)
	}
	if result.ServerPort != 443 {
		t.Errorf("Expected ServerPort=443, got %d", result.ServerPort)
	}
	if result.ServerName != "example.com" {
		t.Errorf("Expected ServerName=example.com, got %s", result.ServerName)
	}
}

// TestUnmarshalShadowsocksConfig 测试 Shadowsocks 反序列化
func TestUnmarshalShadowsocksConfig(t *testing.T) {
	input := []byte(`{
		"data": {
			"id": 1,
			"server_port": 8388,
			"method": "aes-256-gcm",
			"network": "tcp,udp"
		},
		"message": "success"
	}`)

	result, err := UnmarshalShadowsocksConfig(input)
	if err != nil {
		t.Fatalf("UnmarshalShadowsocksConfig() unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if result.ID != 1 {
		t.Errorf("Expected ID=1, got %d", result.ID)
	}
	if result.Method != "aes-256-gcm" {
		t.Errorf("Expected Method=aes-256-gcm, got %s", result.Method)
	}
}

// TestUnmarshalVMessConfig 测试 VMess 反序列化
func TestUnmarshalVMessConfig(t *testing.T) {
	input := []byte(`{
		"data": {
			"id": 1,
			"server_port": 8080,
			"tls": 1,
			"network": "tcp"
		},
		"message": "success"
	}`)

	result, err := UnmarshalVMessConfig(input)
	if err != nil {
		t.Fatalf("UnmarshalVMessConfig() unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if result.ID != 1 {
		t.Errorf("Expected ID=1, got %d", result.ID)
	}
	if result.ServerPort != 8080 {
		t.Errorf("Expected ServerPort=8080, got %d", result.ServerPort)
	}
	if result.TLS != 1 {
		t.Errorf("Expected TLS=1, got %d", result.TLS)
	}
}

// TestUnmarshalAnyTLSConfig 测试 AnyTLS 反序列化
func TestUnmarshalAnyTLSConfig(t *testing.T) {
	input := []byte(`{
		"data": {
			"id": 1,
			"server_port": 443,
			"allow_insecure": 0,
			"server_name": "example.com",
			"padding_rules": "random"
		},
		"message": "success"
	}`)

	result, err := UnmarshalAnyTLSConfig(input)
	if err != nil {
		t.Fatalf("UnmarshalAnyTLSConfig() unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if result.ID != 1 {
		t.Errorf("Expected ID=1, got %d", result.ID)
	}
	if result.ServerName != "example.com" {
		t.Errorf("Expected ServerName=example.com, got %s", result.ServerName)
	}
	if result.PaddingRules != "random" {
		t.Errorf("Expected PaddingRules=random, got %s", result.PaddingRules)
	}
}

// TestUnmarshalUsers 测试用户列表反序列化
func TestUnmarshalUsers(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		wantErr  bool
		validate func(*testing.T, *[]User)
	}{
		{
			name: "valid users list",
			input: []byte(`{
				"data": [
					{"id": 1, "uuid": "uuid-1"},
					{"id": 2, "uuid": "uuid-2"},
					{"id": 3, "uuid": "uuid-3"}
				],
				"message": "success"
			}`),
			wantErr: false,
			validate: func(t *testing.T, users *[]User) {
				if users == nil {
					t.Errorf("Expected non-nil users")
					return
				}
				if len(*users) != 3 {
					t.Errorf("Expected 3 users, got %d", len(*users))
					return
				}
				if (*users)[0].ID != 1 || (*users)[0].UUID != "uuid-1" {
					t.Errorf("User 0 data mismatch")
				}
				if (*users)[1].ID != 2 || (*users)[1].UUID != "uuid-2" {
					t.Errorf("User 1 data mismatch")
				}
			},
		},
		{
			name: "empty users list",
			input: []byte(`{
				"data": [],
				"message": "success"
			}`),
			wantErr: false,
			validate: func(t *testing.T, users *[]User) {
				if users == nil {
					t.Errorf("Expected non-nil users")
					return
				}
				if len(*users) != 0 {
					t.Errorf("Expected 0 users, got %d", len(*users))
				}
			},
		},
		{
			name:    "invalid JSON",
			input:   []byte(`{invalid}`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := UnmarshalUsers(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("UnmarshalUsers() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("UnmarshalUsers() unexpected error: %v", err)
				}
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}
		})
	}
}

// TestMarshalTraffics 测试流量数据序列化
func TestMarshalTraffics(t *testing.T) {
	tests := []struct {
		name     string
		input    []*UserTraffic
		wantErr  bool
		validate func(*testing.T, []byte)
	}{
		{
			name: "valid traffics",
			input: []*UserTraffic{
				{UID: 1, Upload: 1024, Download: 2048, Count: 1},
				{UID: 2, Upload: 4096, Download: 8192, Count: 2},
			},
			wantErr: false,
			validate: func(t *testing.T, data []byte) {
				var traffics []UserTraffic
				err := json.Unmarshal(data, &traffics)
				if err != nil {
					t.Errorf("Failed to unmarshal result: %v", err)
					return
				}
				if len(traffics) != 2 {
					t.Errorf("Expected 2 traffics, got %d", len(traffics))
					return
				}
				if traffics[0].UID != 1 {
					t.Errorf("Expected UID=1, got %d", traffics[0].UID)
				}
				if traffics[0].Upload != 1024 {
					t.Errorf("Expected Upload=1024, got %d", traffics[0].Upload)
				}
				if traffics[1].Download != 8192 {
					t.Errorf("Expected Download=8192, got %d", traffics[1].Download)
				}
			},
		},
		{
			name:    "empty traffics",
			input:   []*UserTraffic{},
			wantErr: false,
			validate: func(t *testing.T, data []byte) {
				expected := "[]"
				if string(data) != expected {
					t.Errorf("Expected %s, got %s", expected, string(data))
				}
			},
		},
		{
			name:    "nil traffics",
			input:   nil,
			wantErr: false,
			validate: func(t *testing.T, data []byte) {
				expected := "null"
				if string(data) != expected {
					t.Errorf("Expected %s, got %s", expected, string(data))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := MarshalTraffics(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("MarshalTraffics() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("MarshalTraffics() unexpected error: %v", err)
				}
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}
		})
	}
}

// TestMarshalTrafficStats 测试流量统计序列化
func TestMarshalTrafficStats(t *testing.T) {
	tests := []struct {
		name     string
		input    *TrafficStats
		wantErr  bool
		validate func(*testing.T, []byte)
	}{
		{
			name: "valid traffic stats",
			input: &TrafficStats{
				Count:    100,
				Requests: 50,
				UserIds:  []int{1, 2, 3},
				UserRequests: map[int]int{
					1: 10,
					2: 20,
					3: 20,
				},
			},
			wantErr: false,
			validate: func(t *testing.T, data []byte) {
				var stats TrafficStats
				err := json.Unmarshal(data, &stats)
				if err != nil {
					t.Errorf("Failed to unmarshal result: %v", err)
					return
				}
				if stats.Count != 100 {
					t.Errorf("Expected Count=100, got %d", stats.Count)
				}
				if stats.Requests != 50 {
					t.Errorf("Expected Requests=50, got %d", stats.Requests)
				}
				if len(stats.UserIds) != 3 {
					t.Errorf("Expected 3 UserIds, got %d", len(stats.UserIds))
				}
				if len(stats.UserRequests) != 3 {
					t.Errorf("Expected 3 UserRequests, got %d", len(stats.UserRequests))
				}
				if stats.UserRequests[1] != 10 {
					t.Errorf("Expected UserRequests[1]=10, got %d", stats.UserRequests[1])
				}
			},
		},
		{
			name: "empty traffic stats",
			input: &TrafficStats{
				Count:        0,
				Requests:     0,
				UserIds:      []int{},
				UserRequests: map[int]int{},
			},
			wantErr: false,
			validate: func(t *testing.T, data []byte) {
				var stats TrafficStats
				err := json.Unmarshal(data, &stats)
				if err != nil {
					t.Errorf("Failed to unmarshal result: %v", err)
					return
				}
				if stats.Count != 0 {
					t.Errorf("Expected Count=0, got %d", stats.Count)
				}
			},
		},
		{
			name:    "nil traffic stats",
			input:   nil,
			wantErr: false,
			validate: func(t *testing.T, data []byte) {
				expected := "null"
				if string(data) != expected {
					t.Errorf("Expected %s, got %s", expected, string(data))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := MarshalTrafficStats(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("MarshalTrafficStats() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("MarshalTrafficStats() unexpected error: %v", err)
				}
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}
		})
	}
}

// TestNodeConfigInterface 测试 NodeConfig 接口的实现
func TestNodeConfigInterface(t *testing.T) {
	tests := []struct {
		name       string
		config     NodeConfig
		wantType   string
		wantString bool // 是否检查 String() 方法返回非空
	}{
		{
			name: "VMessConfig",
			config: &VMessConfig{
				ID:         1,
				ServerPort: 8080,
			},
			wantType:   "vmess",
			wantString: true,
		},
		{
			name: "Hysteria2Config",
			config: &Hysteria2Config{
				ID:         1,
				ServerPort: 8080,
			},
			wantType:   "hysteria2",
			wantString: true,
		},
		{
			name: "HysteriaConfig",
			config: &HysteriaConfig{
				ID:         1,
				ServerPort: 8080,
			},
			wantType:   "hysteria",
			wantString: true,
		},
		{
			name: "TrojanConfig",
			config: &TrojanConfig{
				ID:         1,
				ServerPort: 443,
			},
			wantType:   "trojan",
			wantString: true,
		},
		{
			name: "ShadowsocksConfig",
			config: &ShadowsocksConfig{
				ID:         1,
				ServerPort: 8388,
			},
			wantType:   "shadowsocks",
			wantString: true,
		},
		{
			name: "AnyTLSConfig",
			config: &AnyTLSConfig{
				ID:         1,
				ServerPort: 443,
			},
			wantType:   "anytls",
			wantString: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typeName := tt.config.TypeName()
			if typeName != tt.wantType {
				t.Errorf("Expected TypeName()=%s, got %s", tt.wantType, typeName)
			}

			if tt.wantString {
				str := tt.config.String()
				if str == "" {
					t.Errorf("Expected non-empty String()")
				}
			}
		})
	}
}
