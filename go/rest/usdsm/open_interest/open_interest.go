// Package openinterest implements the Binance USDⓈ-M Futures open-interest endpoint.
package openinterest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/k4k3ru-hub/binance/go/rest/endpoint"
	"github.com/k4k3ru-hub/binance/go/rest/transport"
)

type Client struct{ executor transport.Executor }

type Params struct {
	Symbol string
}

type OpenInterest struct {
	OpenInterest string `json:"openInterest"`
	Symbol       string `json:"symbol"`
	Time         int64  `json:"time"`
}

// NewClient creates a USDⓈ-M Futures open-interest client.
//
// Parameters:
//   - executor: REST request executor.
//
// Returns:
//   - Open-interest client.
//
// Version:
//   - 2026-08-19: Added.
func NewClient(executor transport.Executor) (*Client, error) {
	if executor == nil {
		return nil, fmt.Errorf("failed to create USDⓈ-M open interest client: executor=null")
	}
	return &Client{executor: executor}, nil
}

// Send requests the current open interest for a USDⓈ-M Futures symbol.
//
// Parameters:
//   - ctx: Context for the operation.
//   - params: Request parameters.
//
// Returns:
//   - Current open-interest observation.
//
// Version:
//   - 2026-08-19: Added.
func (c *Client) Send(ctx context.Context, params Params) (*OpenInterest, error) {
	if c == nil || c.executor == nil {
		return nil, fmt.Errorf("failed to request USDⓈ-M open interest: client=null")
	}
	symbol := strings.TrimSpace(params.Symbol)
	if symbol == "" {
		return nil, fmt.Errorf("failed to request USDⓈ-M open interest: symbol=empty")
	}

	body, err := c.executor.Do(ctx, transport.Request{
		Method: http.MethodGet,
		Path:   endpoint.USDSMOpenInterestPath,
		Query:  url.Values{"symbol": {symbol}},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to request USDⓈ-M open interest: %w", err)
	}
	var result OpenInterest
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to request USDⓈ-M open interest: failed to decode response body: %w", err)
	}
	return &result, nil
}
