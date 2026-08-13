// Package rest provides the composed Binance REST client.
package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/k4k3ru-hub/binance/go/rest/endpoint"
	"github.com/k4k3ru-hub/binance/go/rest/spot"
	"github.com/k4k3ru-hub/binance/go/rest/transport"
	"github.com/k4k3ru-hub/binance/go/rest/usdsm"
)

// HTTPClient executes HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// ClientOption configures a REST client.
type ClientOption struct {
	SpotBaseURL    string
	USDSMBaseURL   string
	HTTPClient     HTTPClient
	ConnectTimeout time.Duration
}

// Client is the composition root for Binance REST APIs.
type Client struct {
	spot  *spot.Client
	usdsm *usdsm.Client
}

// DefaultClientOption returns the default REST client options.
func DefaultClientOption() *ClientOption {
	return &ClientOption{
		SpotBaseURL:    endpoint.DefaultSpotBaseURL,
		USDSMBaseURL:   endpoint.DefaultUSDSMBaseURL,
		ConnectTimeout: 3 * time.Second,
	}
}

// NewClient creates a REST client and composes its API groups.
func NewClient(option *ClientOption) (*Client, error) {
	if option == nil {
		option = DefaultClientOption()
	}

	spotBaseURL := strings.TrimRight(option.SpotBaseURL, "/")
	if spotBaseURL == "" {
		spotBaseURL = endpoint.DefaultSpotBaseURL
	}
	usdsmBaseURL := strings.TrimRight(option.USDSMBaseURL, "/")
	if usdsmBaseURL == "" {
		usdsmBaseURL = endpoint.DefaultUSDSMBaseURL
	}

	httpClient := option.HTTPClient
	if httpClient == nil {
		timeout := option.ConnectTimeout
		if timeout == 0 {
			timeout = DefaultClientOption().ConnectTimeout
		}
		if timeout < 0 {
			return nil, fmt.Errorf("failed to create REST client: connect_timeout=out_of_range")
		}
		httpClient = &http.Client{Transport: &http.Transport{DialContext: (&net.Dialer{Timeout: timeout}).DialContext}}
	}

	spotClient, err := spot.NewClient(&executor{baseURL: spotBaseURL, httpClient: httpClient})
	if err != nil {
		return nil, fmt.Errorf("failed to create REST client: %w", err)
	}
	usdsmClient, err := usdsm.NewClient(&executor{baseURL: usdsmBaseURL, httpClient: httpClient})
	if err != nil {
		return nil, fmt.Errorf("failed to create REST client: %w", err)
	}

	return &Client{spot: spotClient, usdsm: usdsmClient}, nil
}

// Spot returns the composed Spot API client.
func (c *Client) Spot() *spot.Client {
	if c == nil {
		return nil
	}
	return c.spot
}

// USDSM returns the composed USDⓈ-M Futures API client.
func (c *Client) USDSM() *usdsm.Client {
	if c == nil {
		return nil
	}
	return c.usdsm
}

type executor struct {
	baseURL    string
	httpClient HTTPClient
}

func (e *executor) Do(ctx context.Context, request transport.Request) ([]byte, error) {
	if e == nil || e.httpClient == nil {
		return nil, fmt.Errorf("failed to execute REST request: executor=null")
	}
	if request.Method == "" || request.Path == "" {
		return nil, fmt.Errorf("failed to execute REST request: invalid request")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	requestURL := e.baseURL + "/" + strings.TrimLeft(request.Path, "/")
	if len(request.Query) != 0 {
		requestURL += "?" + request.Query.Encode()
	}
	httpRequest, err := http.NewRequestWithContext(ctx, request.Method, requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to execute REST request: failed to create request: %w", err)
	}
	httpRequest.Header = request.Header.Clone()

	response, err := e.httpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to execute REST request: %w", err)
	}
	if response == nil || response.Body == nil {
		return nil, fmt.Errorf("failed to execute REST request: response=null")
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to execute REST request: failed to read response body: %w", err)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		responseError := &ResponseError{StatusCode: response.StatusCode}
		_ = json.Unmarshal(body, responseError)
		return nil, responseError
	}
	return body, nil
}

// ResponseError represents a non-successful Binance REST response.
type ResponseError struct {
	StatusCode int    `json:"-"`
	Code       int    `json:"code"`
	Message    string `json:"msg"`
}

func (e *ResponseError) Error() string {
	if e == nil {
		return "failed to execute REST request: response_error=null"
	}
	return fmt.Sprintf("failed to execute REST request: status_code=%d code=%d message=%q", e.StatusCode, e.Code, e.Message)
}
