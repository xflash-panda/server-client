package pkg

import (
	"fmt"
	"strings"

	"github.com/xflash-panda/server-client/pkg/xray"
)

// API is the interface for different panel's api.

const (
	Trojan      NodeType = "trojan"
	ShadowSocks NodeType = "shadowsocks"
	Hysteria    NodeType = "hysteria"
	Hysteria2   NodeType = "hysteria2"
	VMess       NodeType = "vmess"
	AnyTLS      NodeType = "anytls"
)

// ErrorUserNotModified 用户数据未修改错误 (304)
var ErrorUserNotModified = NewNotModifiedError()

type NodeType string

func (n NodeType) String() string {
	return strings.ToLower(string(n))
}

type NodeId int

type configFactoryFunc func() NodeConfig

// 定义一个映射表，将 NodeType 映射到对应的配置工厂函数
var configFactories = map[NodeType]configFactoryFunc{
	Hysteria2:   func() NodeConfig { return &Hysteria2Config{} },
	Hysteria:    func() NodeConfig { return &HysteriaConfig{} },
	Trojan:      func() NodeConfig { return &TrojanConfig{} },
	ShadowSocks: func() NodeConfig { return &ShadowsocksConfig{} },
	VMess:       func() NodeConfig { return &VMessConfig{} },
	AnyTLS:      func() NodeConfig { return &AnyTLSConfig{} },
}

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
	WebSocketConfig *xray.WebSocketConfig `json:"ws_settings,omitempty"`
	GrpcConfig      *xray.GRPCConfig      `json:"grpc_settings,omitempty"`
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
	TlsConfig       *xray.TLSConfig       `json:"tls_settings"`
	WebSocketConfig *xray.WebSocketConfig `json:"ws_settings,omitempty"`
	H2Config        *xray.HTTPConfig      `json:"h2_config"`
	TcpConfig       *xray.TCPConfig       `json:"tcp_settings,omitempty"`
	GrpcConfig      *xray.GRPCConfig      `json:"grpc_settings,omitempty"`
	RouterSettings  *xray.RouterConfig    `json:"router_settings,omitempty"`
	DnsSettings     *xray.DNSConfig       `json:"dns_settings,omitempty"`
}

func (n *VMessConfig) String() string {
	return fmt.Sprintf("VmessConfig: %#v", n)
}

func (n *VMessConfig) TypeName() string {
	return string(VMess)
}

type AnyTLSConfig struct {
	ID            int    `json:"id"`
	ServerPort    int    `json:"server_port"`
	AllowInsecure int    `json:"allow_insecure"`
	ServerName    string `json:"server_name"`
	PaddingRules  string `json:"padding_rules"`
}

func (n *AnyTLSConfig) String() string {
	return fmt.Sprintf("AnyTLSConfig: %#v", n)
}

func (n *AnyTLSConfig) TypeName() string {
	return string(AnyTLS)
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
type (
	RespHeartBeat            RespSubmit
	RespSubmitWithAgent      RespSubmit
	RespSubmitStatsWithAgent RespSubmit
)
