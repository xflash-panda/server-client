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
	var zero T

	val := reflect.ValueOf(nc)

	if val.Kind() != reflect.Ptr || val.IsNil() {
		return zero, fmt.Errorf("AsConfig requires a non-nil pointer to a NodeConfig, got %T", nc)
	}

	targetType := reflect.TypeOf(zero)
	if !val.Elem().Type().AssignableTo(targetType) {
		return zero, fmt.Errorf("cannot assert NodeConfig of type %T to %v", nc, targetType)
	}

	config := val.Elem().Interface()

	tConfig, ok := config.(T)
	if !ok {
		return zero, fmt.Errorf("unexpected error when asserting NodeConfig to %v", targetType)
	}
	return tConfig, nil
}

func UnmarshalConfig[T NodeConfig](data []byte) (*T, error) {
	var config T
	var resp RespConfig
	resp = RespConfig{
		Data: config,
	}
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}
	return &config, nil
}

func UnmarshalHysteria2Config(data []byte) (*Hysteria2Config, error) {
	return UnmarshalConfig[Hysteria2Config](data)
}

func UnmarshalHysteriaConfig(data []byte) (*HysteriaConfig, error) {
	return UnmarshalConfig[HysteriaConfig](data)
}

func UnmarshalTrojanConfig(data []byte) (*TrojanConfig, error) {
	return UnmarshalConfig[TrojanConfig](data)
}

func UnmarshalShadowsocksConfig(data []byte) (*ShadowsocksConfig, error) {
	return UnmarshalConfig[ShadowsocksConfig](data)
}

func UnmarshalVMessConfig(data []byte) (*VMessConfig, error) {
	return UnmarshalConfig[VMessConfig](data)
}
