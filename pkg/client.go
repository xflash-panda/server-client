package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	resty "github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

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
	eTags  sync.Map
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
	path := fmt.Sprintf("/api/v1/server/%s/config", nodeType)
	url := c.assembleURL(path)
	res, err := c.client.R().
		ForceContentType("application/json").
		SetQueryParam("node_id", strconv.Itoa(int(nodeId))).
		Get(path)
	if err != nil {
		return nil, NewNetworkError("request failed", url, err)
	}

	if res.StatusCode() >= 400 {
		body := res.Body()
		return nil, NewAPIErrorFromStatusCode(res.StatusCode(), string(body), url, nil)
	}

	return res.Body(), nil
}

func (c *Client) Config(nodeId NodeId, nodeType NodeType) (config NodeConfig, err error) {
	rawData, err := c.RawConfig(nodeId, nodeType)
	if err != nil {
		return nil, err // 错误已经被 RawConfig 格式化，直接返回即可
	}

	factoryFunc, ok := configFactories[NodeType(nodeType.String())]
	if !ok {
		return nil, NewBusinessLogicError(fmt.Sprintf("invalid config type: %s", nodeType), "")
	}

	var resp RespConfig = RespConfig{
		Data: factoryFunc(),
	}

	if err := json.Unmarshal(rawData, &resp); err != nil {
		return nil, NewParseError("parse response failed", err)
	}

	return resp.Data, nil
}

func (c *Client) RawUsers(nodeId NodeId, nodeType NodeType) (rawData []byte, hash string, err error) {
	path := fmt.Sprintf("/api/v1/server/%s/users", nodeType)
	url := c.assembleURL(path)
	eTagKey := fmt.Sprintf("users_%s_%d", nodeType, nodeId)
	var eTagValue string
	if value, ok := c.eTags.Load(eTagKey); ok {
		eTagValue = value.(string)
	}
	res, err := c.client.R().SetQueryParam("node_id", strconv.Itoa(int(nodeId))).SetHeader("If-None-Match", eTagValue).ForceContentType("application/json").Get(path)
	if err != nil {
		return nil, "", NewNetworkError("request failed", url, err)
	}

	if res.StatusCode() == 304 {
		return nil, "", ErrorUserNotModified
	}

	if res.StatusCode() >= 400 {
		body := res.Body()
		return nil, "", NewAPIErrorFromStatusCode(res.StatusCode(), string(body), url, nil)
	}
	// update etag
	hash = res.Header().Get("Etag")
	c.eTags.Store(eTagKey, hash)
	return res.Body(), hash, nil
}

// Users will pull users form server
func (c *Client) Users(nodeId NodeId, nodeType NodeType) (UserList *[]User, hash string, err error) {
	rawData, hash, err := c.RawUsers(nodeId, nodeType)
	if err != nil {
		return nil, hash, err
	}
	var resp RespUsers
	if err := json.Unmarshal(rawData, &resp); err != nil {
		return nil, hash, NewParseError("parse response failed", err)
	}

	return resp.Data, hash, nil
}

// Submit reports the user traffic
func (c *Client) Submit(nodeId NodeId, nodeType NodeType, userTraffic []*UserTraffic) error {
	path := fmt.Sprintf("/api/v1/server/%s/submit", nodeType)
	url := c.assembleURL(path)
	res, err := c.client.R().
		ForceContentType("application/json").
		SetQueryParam("node_id", strconv.Itoa(int(nodeId))).
		SetBody(userTraffic).
		Post(path)
	if err != nil {
		return NewNetworkError("request failed", url, err)
	}

	if res.StatusCode() >= 400 {
		body := res.Body()
		return NewAPIErrorFromStatusCode(res.StatusCode(), string(body), url, nil)
	}

	var resp RespSubmit
	if err := json.Unmarshal(res.Body(), &resp); err != nil {
		return NewParseError("parse response failed", err)
	}
	return nil
}

func (c *Client) SubmitWithAgent(nodeId NodeId, nodeType NodeType, userTraffic []*UserTraffic) error {
	path := fmt.Sprintf("/api/v1/server/%s/submitWithAgent", nodeType)
	url := c.assembleURL(path)
	res, err := c.client.R().
		ForceContentType("application/json").
		SetQueryParams(map[string]string{"node_id": strconv.Itoa(int(nodeId))}).
		SetBody(userTraffic).
		Post(path)
	if err != nil {
		return NewNetworkError("request failed", url, err)
	}

	if res.StatusCode() >= 400 {
		body := res.Body()
		return NewAPIErrorFromStatusCode(res.StatusCode(), string(body), url, nil)
	}

	var resp RespSubmitWithAgent
	if err := json.Unmarshal(res.Body(), &resp); err != nil {
		return NewParseError("parse response failed", err)
	}
	return nil
}

func (c *Client) SubmitStatsWithAgent(nodeId NodeId, nodeType NodeType, nodeIp string, stats *TrafficStats) error {
	path := fmt.Sprintf("/api/v1/server/%s/submitStatsWithAgent", nodeType)
	url := c.assembleURL(path)
	res, err := c.client.R().
		ForceContentType("application/json").
		SetQueryParams(map[string]string{"node_id": strconv.Itoa(int(nodeId)), "node_ip": nodeIp}).
		SetBody(stats).
		Post(path)
	if err != nil {
		return NewNetworkError("request failed", url, err)
	}

	if res.StatusCode() >= 400 {
		body := res.Body()
		return NewAPIErrorFromStatusCode(res.StatusCode(), string(body), url, nil)
	}

	var resp RespSubmitWithAgent
	if err := json.Unmarshal(res.Body(), &resp); err != nil {
		return NewParseError("parse response failed", err)
	}
	return nil
}

func (c *Client) Heartbeat(nodeId NodeId, nodeType NodeType, nodeIp string) error {
	path := fmt.Sprintf("/api/v1/server/%s/heartbeat", nodeType)
	url := c.assembleURL(path)
	req := c.client.R().
		ForceContentType("application/json").
		SetQueryParam("node_id", strconv.Itoa(int(nodeId)))

	// 只在 nodeIp 不为空时才添加查询参数
	if nodeIp != "" {
		req.SetQueryParam("node_ip", nodeIp)
	}

	res, err := req.Get(path)
	if err != nil {
		return NewNetworkError("request failed", url, err)
	}

	if res.StatusCode() >= 400 {
		body := res.Body()
		return NewAPIErrorFromStatusCode(res.StatusCode(), string(body), url, nil)
	}

	var respHeartBeat RespHeartBeat
	if err := json.Unmarshal(res.Body(), &respHeartBeat); err != nil {
		return NewParseError("parse response failed", err)
	}
	return nil
}
