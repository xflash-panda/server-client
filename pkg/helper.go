package pkg

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func AsVMessConfig(nc NodeConfig) (*VMessConfig, error) {
	config, err := AsConfig[*VMessConfig](nc)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func AsHysteriaConfig(nc NodeConfig) (*HysteriaConfig, error) {
	config, err := AsConfig[*HysteriaConfig](nc)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func AsHysteria2Config(nc NodeConfig) (*Hysteria2Config, error) {
	config, err := AsConfig[*Hysteria2Config](nc)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func AsTrojanConfig(nc NodeConfig) (*TrojanConfig, error) {
	config, err := AsConfig[*TrojanConfig](nc)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func AsShadowsocksConfig(nc NodeConfig) (*ShadowsocksConfig, error) {
	config, err := AsConfig[*ShadowsocksConfig](nc)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func AsConfig[T NodeConfig](nc NodeConfig) (T, error) {
	// 创建类型 T 的零值
	var zero T

	// 使用反射来获取 nc 的实际类型
	val := reflect.ValueOf(nc)
	if val.Kind() == reflect.Ptr && val.IsNil() {
		// 如果 nc 是 nil 指针，则返回零值和错误
		return zero, fmt.Errorf("nil cannot be converted to type %v", reflect.TypeOf(zero))
	}

	// 使用类型断言尝试将 nc 转换为具体的类型 T
	tConfig, ok := nc.(T)
	if !ok {
		// 如果断言失败，返回零值和错误
		return zero, fmt.Errorf("cannot assert type %v to type %v", reflect.TypeOf(nc), reflect.TypeOf(zero))
	}

	// 如果断言成功，返回结果
	return tConfig, nil
}

func UnmarshalHysteria2Config(data []byte) (*Hysteria2Config, error) {
	var resp RespConfig
	resp = RespConfig{
		Data: &Hysteria2Config{},
	}
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}
	return resp.Data.(*Hysteria2Config), nil
}

func UnmarshalHysteriaConfig(data []byte) (*HysteriaConfig, error) {
	var resp RespConfig
	resp = RespConfig{
		Data: &HysteriaConfig{},
	}
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}
	return resp.Data.(*HysteriaConfig), nil
}

func UnmarshalTrojanConfig(data []byte) (*TrojanConfig, error) {
	var resp RespConfig
	resp = RespConfig{
		Data: &TrojanConfig{},
	}
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}
	return resp.Data.(*TrojanConfig), nil
}

func UnmarshalShadowsocksConfig(data []byte) (*ShadowsocksConfig, error) {
	var resp RespConfig
	resp = RespConfig{
		Data: &ShadowsocksConfig{},
	}
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}
	return resp.Data.(*ShadowsocksConfig), nil
}

func UnmarshalVMessConfig(data []byte) (*VMessConfig, error) {
	var resp RespConfig
	resp = RespConfig{
		Data: &VMessConfig{},
	}
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}
	return resp.Data.(*VMessConfig), nil
}

func UnmarshalUsers(data []byte) (*[]User, error) {
	var resp RespUsers
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}
	return resp.Data, nil
}
