package pkg

import (
	"fmt"
	"github.com/xtls/xray-core/infra/conf"
)

// API is the interface for different panel's api.

type NodeConfig interface {
	String() string
	TypeName() string
}

type UserTraffic struct {
	UID      int    `json:"user_id"`
	Upload   uint64 `json:"u"`
	Download uint64 `json:"d"`
	Count    uint64 `json:"n"`
}

type Hysteria2Config struct {
	ID                 int    `json:"id"`
	ServerPort         int    `json:"server_port"`
	Obfs               string `json:"obfs"`
	UpMbps             int    `json:"up_mbps"`
	DownMbps           int    `json:"down_mbps"`
	IgnoreCliBandWidth bool   `json:"ignore_cli_band_width"`
	DisableUDP         bool   `json:"disable_udp"`
}

func (n *Hysteria2Config) String() string {
	return fmt.Sprintf("Hysteria2Config: %#v", n)
}

func (n *Hysteria2Config) TypeName() string {
	return string(Hysteria2)
}

type HysteriaConfig struct {
	ID                  int    `json:"id"`
	ServerPort          int    `json:"server_port"`
	Protocol            string `json:"protocol"`
	Obfs                string `json:"obfs"`
	UpMbps              int    `json:"up_mbps"`
	DownMbps            int    `json:"down_mbps"`
	DisableMTUDiscovery bool   `json:"disable_mtu_discovery"`
	DisableUdp          bool   `json:"disable_udp"`
}

func (n *HysteriaConfig) TypeName() string {
	return string(Hysteria)
}

func (n *HysteriaConfig) String() string {
	return fmt.Sprintf("HysteriaConfig: %#v", n)
}

type ShadowsocksConfig struct {
	ID         int    `json:"id"`
	ServerPort int    `json:"server_port"`
	Method     string `json:"method"`
	Network    string `json:"network"`
}

func (n *ShadowsocksConfig) String() string {
	return fmt.Sprintf("ShadowSocksConfig: %#v", n)
}

func (n *ShadowsocksConfig) TypeName() string {
	return string(ShadowSocks)
}

type TrojanConfig struct {
	ID              int                   `json:"id"`
	ServerPort      int                   `json:"server_port"`
	AllowInsecure   int                   `json:"allow_insecure"`
	ServerName      string                `json:"server_name"`
	Network         string                `json:"network"`
	WebSocketConfig *conf.WebSocketConfig `json:"ws_settings,omitempty"`
	GrpcConfig      *conf.GRPCConfig      `json:"grpc_settings,omitempty"`
}

func (n *TrojanConfig) String() string {
	return fmt.Sprintf("TrojanConfig: %#v", n)
}

func (n *TrojanConfig) TypeName() string {
	return string(Trojan)
}

type VMessConfig struct {
	ID              int                   `json:"id"`
	ServerPort      int                   `json:"server_port"`
	TLS             int                   `json:"tls"`
	Network         string                `json:"network"`
	TlsConfig       *conf.TLSConfig       `json:"tls_settings"`
	WebSocketConfig *conf.WebSocketConfig `json:"ws_settings,omitempty"`
	H2Config        *conf.HTTPConfig      `json:"h2_config"`
	TcpConfig       *conf.TCPConfig       `json:"tcp_settings,omitempty"`
	GrpcConfig      *conf.GRPCConfig      `json:"grpc_settings,omitempty"`
	RouterSettings  *conf.RouterConfig    `json:"router_settings,omitempty"`
	DnsSettings     *conf.DNSConfig       `json:"dns_settings,omitempty"`
}

func (n *VMessConfig) String() string {
	return fmt.Sprintf("VmessConfig: %#v", n)
}

func (n *VMessConfig) TypeName() string {
	return string(VMess)
}

type User struct {
	ID   int    `json:"id"`
	UUID string `json:"uuid"`
}

type TrafficStats struct {
	Count        int         `json:"count"`
	Requests     int         `json:"requests"`
	UserIds      []int       `json:"user_ids"`
	UserRequests map[int]int `json:"user_requests"`
}

type RespUsers struct {
	Data    *[]User `json:"data"`
	Message string  `json:"message"`
}

type RespConfig struct {
	Data    NodeConfig `json:"data"`
	Message string     `json:"message"`
}

type RespSubmit struct {
	Data    bool   `json:"data"`
	Message string `json:"message"`
}
type RespHeartBeat RespSubmit
type RespSubmitWithAgent RespSubmit
type RespSubmitStatsWithAgent RespSubmit
