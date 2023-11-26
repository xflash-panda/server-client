package pkg

import "C"
import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

const (
	Trojan      NodeType = "trojan"
	ShadowSocks NodeType = "shadowsocks"
	Hysteria    NodeType = "hysteria"
	Hysteria2   NodeType = "hysteria2"
	VMess       NodeType = "vmess"
)

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
}

// Config  api config
type Config struct {
	APIHost string
	Token   string
	Timeout time.Duration
	Debug   bool
}

// Client APIClient create a api client to the panel.
type Client struct {
	client *resty.Client
	config *Config
}

// New creat a api instance
func New(apiConfig *Config) *Client {
	client := resty.New()
	if apiConfig.Timeout > 0 {
		client.SetTimeout(apiConfig.Timeout)
	} else {
		client.SetTimeout(5 * time.Second)
	}
	client.OnError(func(req *resty.Request, err error) {
		var v *resty.ResponseError
		if errors.As(err, &v) {
			// v.Response contains the last response from the server
			// v.Err contains the original error
			log.Errorln(v.Err)
		}
	})
	client.SetBaseURL(apiConfig.APIHost)
	// Create Key for each requests
	client.SetRetryCount(3)
	client.SetQueryParams(map[string]string{
		"token": apiConfig.Token,
	})
	client.SetCloseConnection(true)

	if apiConfig.Debug {
		client.SetDebug(true)
	}

	apiClient := &Client{
		client: client,
		config: apiConfig,
	}
	return apiClient
}

// Debug set the client debug for client
func (c *Client) Debug(enable bool) {
	c.client.SetDebug(enable)
}

func (c *Client) assembleURL(path string) string {
	return c.config.APIHost + path
}

func (c *Client) RawConfig(nodeId NodeId, nodeType NodeType) (rawData []byte, err error) {
	var path = fmt.Sprintf("/api/v1/server/%s/config", nodeType)
	res, err := c.client.R().
		ForceContentType("application/json").
		SetQueryParam("node_id", strconv.Itoa(int(nodeId))).
		Get(path)

	if err != nil {
		return nil, fmt.Errorf("request %s failed: %w", c.assembleURL(path), err)
	}

	if res.StatusCode() >= 400 {
		body := res.Body()
		return nil, fmt.Errorf("request %s failed: %s, %w", c.assembleURL(path), string(body), err)
	}

	return res.Body(), nil
}

func (c *Client) Config(nodeId NodeId, nodeType NodeType) (config NodeConfig, err error) {
	rawData, err := c.RawConfig(nodeId, nodeType)
	if err != nil {
		return nil, err // 错误已经被 GetRawConfig 格式化，直接返回即可
	}

	factoryFunc, ok := configFactories[NodeType(nodeType.String())]
	if !ok {
		return nil, fmt.Errorf("invalid config type: %s", nodeType)
	}

	var resp RespConfig
	resp = RespConfig{
		Data: factoryFunc(),
	}

	if err := json.Unmarshal(rawData, &resp); err != nil {
		return nil, fmt.Errorf("parse response failed: %w", err)
	}

	if len(resp.Message) > 0 {
		return nil, fmt.Errorf("api error, message: %s", resp.Message)
	}

	return resp.Data, nil
}

// Users will pull users form server
func (c *Client) Users(nodeId NodeId, nodeType NodeType) (UserList *[]User, err error) {
	var path = fmt.Sprintf("/api/v1/server/%s/users", nodeType)
	res, err := c.client.R().SetQueryParam("node_id", strconv.Itoa(int(nodeId))).ForceContentType("application/json").Get(path)

	if err != nil {
		return nil, fmt.Errorf("request %s failed: %s", c.assembleURL(path), err)
	}

	if res.StatusCode() > 400 {
		body := res.Body()
		return nil, fmt.Errorf("request %s failed: %s, %s", c.assembleURL(path), string(body), err)
	}
	var resp RespUsers
	if err := json.Unmarshal(res.Body(), &resp); err != nil {
		return nil, fmt.Errorf("parse response failed: %s", err)
	}

	if len(resp.Message) > 0 {
		return nil, fmt.Errorf("api error, message: %s", resp.Message)
	}

	return resp.Data, nil
}

// Submit reports the user traffic
func (c *Client) Submit(nodeId NodeId, nodeType NodeType, userTraffic []*UserTraffic) error {
	var path = fmt.Sprintf("/api/v1/server/%s/submit", nodeType)
	res, err := c.client.R().SetQueryParam("node_id", strconv.Itoa(int(nodeId))).SetBody(userTraffic).Post(path)
	if err != nil {
		return fmt.Errorf("request %s failed: %s", c.assembleURL(path), err)
	}

	if res.StatusCode() > 400 {
		body := res.Body()
		return fmt.Errorf("request %s failed: %s, %s", c.assembleURL(path), string(body), err)
	}

	var resp RespSubmit
	if err := json.Unmarshal(res.Body(), &resp); err != nil {
		return fmt.Errorf("parse response failed: %s", err)
	}
	if len(resp.Message) > 0 {
		return fmt.Errorf("api error, message: %s", resp.Message)
	}
	return nil
}

func (c *Client) SubmitWithAgent(nodeId NodeId, nodeType NodeType, userTraffic []*UserTraffic) error {
	var path = fmt.Sprintf("/api/v1/server/%s/submitWithAgent", nodeType)
	res, err := c.client.R().SetQueryParams(map[string]string{"node_id": strconv.Itoa(int(nodeId))}).SetBody(userTraffic).Post(path)
	if err != nil {
		return fmt.Errorf("request %s failed: %s", c.assembleURL(path), err)
	}

	if res.StatusCode() > 400 {
		body := res.Body()
		return fmt.Errorf("request %s failed: %s, %s", c.assembleURL(path), string(body), err)
	}

	var resp RespSubmitWithAgent
	if err := json.Unmarshal(res.Body(), &resp); err != nil {
		return fmt.Errorf("parse response failed: %s", err)
	}
	if len(resp.Message) > 0 {
		return fmt.Errorf("api error, message: %s", resp.Message)
	}
	return nil
}

func (c *Client) SubmitStatsWithAgent(nodeId NodeId, nodeType NodeType, nodeIp string, stats *TrafficStats) error {
	var path = fmt.Sprintf("/api/v1/server/%s/submitStatsWithAgent", nodeType)
	res, err := c.client.R().SetQueryParams(map[string]string{"node_id": strconv.Itoa(int(nodeId)), "node_ip": nodeIp}).SetBody(stats).Post(path)
	if err != nil {
		return fmt.Errorf("request %s failed: %s", c.assembleURL(path), err)
	}

	if res.StatusCode() > 400 {
		body := res.Body()
		return fmt.Errorf("request %s failed: %s, %s", c.assembleURL(path), string(body), err)
	}

	var resp RespSubmitWithAgent
	if err := json.Unmarshal(res.Body(), &resp); err != nil {
		return fmt.Errorf("parse response failed: %s", err)
	}
	if len(resp.Message) > 0 {
		return fmt.Errorf("api error, message: %s", resp.Message)
	}
	return nil
}

func (c *Client) Heartbeat(nodeId NodeId, nodeType NodeType, nodeIp string) error {
	var path = fmt.Sprintf("/api/v1/server/%s/heartbeat", nodeType)
	res, err := c.client.R().
		ForceContentType("application/json").
		SetQueryParams(map[string]string{"node_id": strconv.Itoa(int(nodeId)), "node_ip": nodeIp}).
		Get(path)

	if err != nil {
		return fmt.Errorf("request %s failed: %s", c.assembleURL(path), err)
	}

	if res.StatusCode() > 400 {
		body := res.Body()
		return fmt.Errorf("request %s failed: %s, %s", c.assembleURL(path), string(body), err)
	}

	var respHeartBeat RespHeartBeat
	if err := json.Unmarshal(res.Body(), &respHeartBeat); err != nil {
		return fmt.Errorf("parse response failed: %s", err)
	}
	if len(respHeartBeat.Message) > 0 {
		return fmt.Errorf("api error, message: %s", respHeartBeat.Message)
	}
	return nil
}
