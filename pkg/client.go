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

func (c *Client) Register(nodeId NodeId, nodeType NodeType, hostname string, port int, nodeIp string) (registerId int, config NodeConfig, err error) {
	path := fmt.Sprintf("/api/v1/server/enhanced/%s/register", nodeType)
	url := c.assembleURL(path)

	body := map[string]interface{}{"hostname": hostname, "port": port}
	if nodeIp != "" {
		body["node_ip"] = nodeIp
	}

	res, err := c.client.R().
		ForceContentType("application/json").
		SetQueryParam("node_id", strconv.Itoa(int(nodeId))).
		SetBody(body).
		Post(path)
	if err != nil {
		return 0, nil, NewNetworkError("request failed", url, err)
	}

	if res.StatusCode() >= 400 {
		body := res.Body()
		return 0, nil, NewAPIErrorFromStatusCode(res.StatusCode(), string(body), url, nil)
	}

	factoryFunc, ok := configFactories[NodeType(nodeType.String())]
	if !ok {
		return 0, nil, NewBusinessLogicError(fmt.Sprintf("invalid config type: %s", nodeType), "")
	}

	var resp RespRegister = RespRegister{
		Config: factoryFunc(),
	}

	if err := json.Unmarshal(res.Body(), &resp); err != nil {
		return 0, nil, NewParseError("parse response failed", err)
	}

	return resp.RegisterId, resp.Config, nil
}

func (c *Client) Unregister(nodeType NodeType, registerId int) error {
	path := fmt.Sprintf("/api/v1/server/enhanced/%s/unregister", nodeType)
	url := c.assembleURL(path)

	res, err := c.client.R().
		ForceContentType("application/json").
		SetQueryParam("register_id", strconv.Itoa(registerId)).
		Post(path)
	if err != nil {
		return NewNetworkError("request failed", url, err)
	}

	if res.StatusCode() >= 400 {
		body := res.Body()
		return NewAPIErrorFromStatusCode(res.StatusCode(), string(body), url, nil)
	}

	return nil
}

func (c *Client) RawUsers(registerId int, nodeType NodeType) (rawData []byte, err error) {
	path := fmt.Sprintf("/api/v1/server/enhanced/%s/users", nodeType)
	url := c.assembleURL(path)
	eTagKey := fmt.Sprintf("users_%s_%d", nodeType, registerId)
	var eTagValue string
	if value, ok := c.eTags.Load(eTagKey); ok {
		eTagValue = value.(string)
	}
	res, err := c.client.R().SetQueryParam("register_id", strconv.Itoa(registerId)).SetHeader("If-None-Match", eTagValue).ForceContentType("application/json").Get(path)
	if err != nil {
		return nil, NewNetworkError("request failed", url, err)
	}

	if res.StatusCode() == 304 {
		return nil, ErrorUserNotModified
	}

	if res.StatusCode() >= 400 {
		body := res.Body()
		return nil, NewAPIErrorFromStatusCode(res.StatusCode(), string(body), url, nil)
	}
	// update etag
	hash := res.Header().Get("Etag")
	c.eTags.Store(eTagKey, hash)
	return res.Body(), nil
}

// Users will pull users form server
func (c *Client) Users(registerId int, nodeType NodeType) (UserList *[]User, err error) {
	rawData, err := c.RawUsers(registerId, nodeType)
	if err != nil {
		return nil, err
	}
	var resp RespUsers
	if err := json.Unmarshal(rawData, &resp); err != nil {
		return nil, NewParseError("parse response failed", err)
	}

	return resp.Data, nil
}

// Submit reports the user traffic
func (c *Client) Submit(registerId int, nodeType NodeType, userTraffic []*UserTraffic) error {
	path := fmt.Sprintf("/api/v1/server/enhanced/%s/submit", nodeType)
	url := c.assembleURL(path)

	body := map[string]interface{}{
		"register_id": registerId,
		"data":        userTraffic,
	}

	res, err := c.client.R().
		ForceContentType("application/json").
		SetBody(body).
		Post(path)
	if err != nil {
		return NewNetworkError("request failed", url, err)
	}

	if res.StatusCode() >= 400 {
		respBody := res.Body()
		return NewAPIErrorFromStatusCode(res.StatusCode(), string(respBody), url, nil)
	}

	var resp RespSubmit
	if err := json.Unmarshal(res.Body(), &resp); err != nil {
		return NewParseError("parse response failed", err)
	}
	return nil
}

func (c *Client) SubmitWithAgent(registerId int, nodeType NodeType, userTraffic []*UserTraffic) error {
	path := fmt.Sprintf("/api/v1/server/enhanced/%s/submitWithAgent", nodeType)
	url := c.assembleURL(path)

	body := map[string]interface{}{
		"register_id": registerId,
		"data":        userTraffic,
	}

	res, err := c.client.R().
		ForceContentType("application/json").
		SetBody(body).
		Post(path)
	if err != nil {
		return NewNetworkError("request failed", url, err)
	}

	if res.StatusCode() >= 400 {
		respBody := res.Body()
		return NewAPIErrorFromStatusCode(res.StatusCode(), string(respBody), url, nil)
	}

	var resp RespSubmitWithAgent
	if err := json.Unmarshal(res.Body(), &resp); err != nil {
		return NewParseError("parse response failed", err)
	}
	return nil
}

func (c *Client) SubmitStatsWithAgent(registerId int, nodeType NodeType, nodeIp string, stats *TrafficStats) error {
	path := fmt.Sprintf("/api/v1/server/enhanced/%s/submitStatsWithAgent", nodeType)
	url := c.assembleURL(path)

	body := map[string]interface{}{
		"register_id": registerId,
		"data":        stats,
	}
	if nodeIp != "" {
		body["node_ip"] = nodeIp
	}

	res, err := c.client.R().
		ForceContentType("application/json").
		SetBody(body).
		Post(path)
	if err != nil {
		return NewNetworkError("request failed", url, err)
	}

	if res.StatusCode() >= 400 {
		respBody := res.Body()
		return NewAPIErrorFromStatusCode(res.StatusCode(), string(respBody), url, nil)
	}

	var resp RespSubmitWithAgent
	if err := json.Unmarshal(res.Body(), &resp); err != nil {
		return NewParseError("parse response failed", err)
	}
	return nil
}

func (c *Client) Heartbeat(registerId int, nodeType NodeType, nodeIp string) error {
	path := fmt.Sprintf("/api/v1/server/enhanced/%s/heartbeat", nodeType)
	url := c.assembleURL(path)

	body := map[string]interface{}{"register_id": registerId}
	if nodeIp != "" {
		body["node_ip"] = nodeIp
	}

	res, err := c.client.R().
		ForceContentType("application/json").
		SetBody(body).
		Post(path)
	if err != nil {
		return NewNetworkError("request failed", url, err)
	}

	if res.StatusCode() >= 400 {
		respBody := res.Body()
		return NewAPIErrorFromStatusCode(res.StatusCode(), string(respBody), url, nil)
	}

	var respHeartBeat RespHeartBeat
	if err := json.Unmarshal(res.Body(), &respHeartBeat); err != nil {
		return NewParseError("parse response failed", err)
	}
	return nil
}
