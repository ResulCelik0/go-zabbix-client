// Zabbix API client for Go
//
// Author: Resul Çelik 2024

package gozabbix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Request struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Id      int         `json:"id"`
}

type Response struct {
	Jsonrpc string       `json:"jsonrpc"`
	Result  interface{}  `json:"result"`
	Error   *ZabbixError `json:"error,omitempty"`
	Id      int          `json:"id"`
}

type ZabbixError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func (e *ZabbixError) Error() string {
	return fmt.Sprintf("zabbix returned error; code: %d, message: %s, data: %s", e.Code, e.Message, e.Data)
}

type Config struct {
	URL      string
	Username string
	Password string
}

type ZabbixClient struct {
	token        string
	zabbixConfig *Config
	httpClient   *http.Client
}

// NewZabbixClient creates a new Zabbix client with the given configuration.
// If getUserInfo is true, the client will also fetch the user information. userInfo will be nil if getUserInfo is false.
func NewZabbixClient(config *Config, getUserInfo bool) (client *ZabbixClient, userInfo *LoginResponse, err error) {
	client = &ZabbixClient{
		zabbixConfig: config,
		httpClient:   &http.Client{},
	}
	token, userInfo, err := client.UserAPI().Login(&LoginRequest{
		Username: config.Username,
		Password: config.Password,
		UserData: getUserInfo,
	})
	if err != nil {
		return
	}
	client.token = token
	return
}

func (s *ZabbixClient) execute(req *Request) (resp *Response, err error) {
	req.Jsonrpc = "2.0"
	if req.Id == 0 {
		req.Id = 1
	}
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	httpRequest, err := http.NewRequest("POST", s.zabbixConfig.URL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("User-Agent", "gozabbix")
	httpRequest.Header.Set("Accept", "application/json")
	if s.token != "" {
		httpRequest.Header.Set("Authorization", "Bearer "+s.token)
	}
	httResponse, err := s.httpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer httResponse.Body.Close()
	rawBody, err := io.ReadAll(httResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	if httResponse.StatusCode != 200 {
		return nil, fmt.Errorf("http status code: %d, body: %s", httResponse.StatusCode, string(rawBody))
	}
	err = json.Unmarshal(rawBody, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	return
}

func (s *ZabbixClient) UserAPI() *UserAPI {
	return &UserAPI{s}
}
