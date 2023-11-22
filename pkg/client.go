package api

import "C"
import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type NodeType string

const (
	Trojan      NodeType = "trojan"
	ShadowSocks NodeType = "shadowsocks"
	Hysteria    NodeType = "hysteria"
	Hysteria2   NodeType = "hysteria2"
	VMess       NodeType = "vmess"
)

// Config  api config
type Config struct {
	APIHost  string
	NodeID   int
	NodeType NodeType
	Token    string
	Timeout  time.Duration
	Debug    bool
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
		"node_id": strconv.Itoa(apiConfig.NodeID),
		"token":   apiConfig.Token,
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

// Config will pull config form server
func (c *Client) Config() (config NodeConfig, err error) {
	var path = fmt.Sprintf("/api/v1/server/%s/config", c.config.NodeType)
	res, err := c.client.R().
		ForceContentType("application/json").
		Get(path)

	if err != nil {
		return nil, fmt.Errorf("request %s failed: %s", c.assembleURL(path), err)
	}

	if res.StatusCode() > 400 {
		body := res.Body()
		return nil, fmt.Errorf("request %s failed: %s, %s", c.assembleURL(path), string(body), err)
	}

	var resp RespConfig
	switch c.config.NodeType {
	case Hysteria2:
		resp.Data = &Hysteria2Config{}
		break
	}
	if err := json.Unmarshal(res.Body(), &resp); err != nil {
		return nil, fmt.Errorf("parse response failed: %s", err)
	}

	if len(resp.Message) > 0 {
		return nil, fmt.Errorf("api error, message: %s", resp.Message)
	}
	return resp.Data, nil
}

// Users will pull users form server
func (c *Client) Users() (UserList *[]User, err error) {
	var path = fmt.Sprintf("/api/v1/server/%s/users", c.config.NodeType)
	res, err := c.client.R().ForceContentType("application/json").Get(path)

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
func (c *Client) Submit(userTraffic []*UserTraffic) error {
	var path = fmt.Sprintf("/api/v1/server/%s/submit", c.config.NodeType)
	res, err := c.client.R().SetBody(userTraffic).Post(path)
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

func (c *Client) SubmitWithAgent(nodeIp string, userTraffic []*UserTraffic) error {
	var path = fmt.Sprintf("/api/v1/server/%s/submitWithAgent", c.config.NodeType)
	res, err := c.client.R().SetQueryParam("node_ip", nodeIp).SetBody(userTraffic).Post(path)
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

func (c *Client) SubmitStatsWithAgent(nodeIp string, stats *Stats) error {
	var path = fmt.Sprintf("/api/v1/server/%s/submitStatsWithAgent", c.config.NodeType)
	res, err := c.client.R().SetQueryParam("node_ip", nodeIp).SetBody(stats).Post(path)
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

func (c *Client) Heartbeat(nodeIp string) error {
	var path = fmt.Sprintf("/api/v1/server/%s/heartbeat", c.config.NodeType)
	res, err := c.client.R().
		ForceContentType("application/json").
		SetQueryParam("node_ip", nodeIp).
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
