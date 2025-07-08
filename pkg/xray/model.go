package xray

import (
	"encoding/json"
)

type StringList []string

type TLSConfig struct {
	Insecure                             bool             `json:"allowInsecure"`
	Certs                                []*TLSCertConfig `json:"certificates"`
	ServerName                           string           `json:"serverName"`
	ALPN                                 *StringList      `json:"alpn"`
	EnableSessionResumption              bool             `json:"enableSessionResumption"`
	DisableSystemRoot                    bool             `json:"disableSystemRoot"`
	MinVersion                           string           `json:"minVersion"`
	MaxVersion                           string           `json:"maxVersion"`
	CipherSuites                         string           `json:"cipherSuites"`
	PreferServerCipherSuites             bool             `json:"preferServerCipherSuites"`
	Fingerprint                          string           `json:"fingerprint"`
	RejectUnknownSNI                     bool             `json:"rejectUnknownSni"`
	PinnedPeerCertificateChainSha256     *[]string        `json:"pinnedPeerCertificateChainSha256"`
	PinnedPeerCertificatePublicKeySha256 *[]string        `json:"pinnedPeerCertificatePublicKeySha256"`
}

type WebSocketConfig struct {
	Path                string            `json:"path"`
	Headers             map[string]string `json:"headers"`
	AcceptProxyProtocol bool              `json:"acceptProxyProtocol"`
}

type HTTPConfig struct {
	Host               *StringList            `json:"host"`
	Path               string                 `json:"path"`
	ReadIdleTimeout    int32                  `json:"read_idle_timeout"`
	HealthCheckTimeout int32                  `json:"health_check_timeout"`
	Method             string                 `json:"method"`
	Headers            map[string]*StringList `json:"headers"`
}

type TLSCertConfig struct {
	CertFile       string   `json:"certificateFile"`
	CertStr        []string `json:"certificate"`
	KeyFile        string   `json:"keyFile"`
	KeyStr         []string `json:"key"`
	Usage          string   `json:"usage"`
	OcspStapling   uint64   `json:"ocspStapling"`
	OneTimeLoading bool     `json:"oneTimeLoading"`
}

type TCPConfig struct {
	HeaderConfig        json.RawMessage `json:"header"`
	AcceptProxyProtocol bool            `json:"acceptProxyProtocol"`
}

type GRPCConfig struct {
	ServiceName         string `json:"serviceName" `
	MultiMode           bool   `json:"multiMode"`
	IdleTimeout         int32  `json:"idle_timeout"`
	HealthCheckTimeout  int32  `json:"health_check_timeout"`
	PermitWithoutStream bool   `json:"permit_without_stream"`
	InitialWindowsSize  int32  `json:"initial_windows_size"`
	UserAgent           string `json:"user_agent"`
}

type RouterConfig struct {
	Settings       *RouterRulesConfig `json:"settings"` // Deprecated
	RuleList       []json.RawMessage  `json:"rules"`
	DomainStrategy *string            `json:"domainStrategy"`
	Balancers      []*BalancingRule   `json:"balancers"`

	DomainMatcher string `json:"domainMatcher"`
}

type RouterRulesConfig struct {
	RuleList       []json.RawMessage `json:"rules"`
	DomainStrategy string            `json:"domainStrategy"`
}

type BalancingRule struct {
	Tag       string         `json:"tag"`
	Selectors StringList     `json:"selector"`
	Strategy  StrategyConfig `json:"strategy"`
}

// StrategyConfig represents a strategy config
type StrategyConfig struct {
	Type     string           `json:"type"`
	Settings *json.RawMessage `json:"settings"`
}

type DNSConfig struct {
	Servers                []*NameServerConfig `json:"servers"`
	Hosts                  *HostsWrapper       `json:"hosts"`
	ClientIP               *Address            `json:"clientIp"`
	Tag                    string              `json:"tag"`
	QueryStrategy          string              `json:"queryStrategy"`
	DisableCache           bool                `json:"disableCache"`
	DisableFallback        bool                `json:"disableFallback"`
	DisableFallbackIfMatch bool                `json:"disableFallbackIfMatch"`
}

type NameServerConfig struct {
	Address       *Address
	ClientIP      *Address
	Port          uint16
	SkipFallback  bool
	Domains       []string
	ExpectIPs     StringList
	QueryStrategy string
}

type HostsWrapper struct {
	Hosts map[string]*HostAddress
}

type Address struct{}

type HostAddress struct {
	_ *Address
	_ []*Address
}
