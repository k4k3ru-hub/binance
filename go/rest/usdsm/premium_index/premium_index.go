// Package premiumindex implements the Binance USDⓈ-M Futures mark-price and funding-rate endpoint.
package premiumindex

import (
	"bytes"
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

type PremiumIndex struct {
	Symbol               string `json:"symbol"`
	MarkPrice            string `json:"markPrice"`
	IndexPrice           string `json:"indexPrice"`
	EstimatedSettlePrice string `json:"estimatedSettlePrice"`
	LastFundingRate      string `json:"lastFundingRate"`
	InterestRate         string `json:"interestRate"`
	NextFundingTime      int64  `json:"nextFundingTime"`
	Time                 int64  `json:"time"`
}

// NewClient creates a USDⓈ-M Futures premium-index client.
//
// Parameters:
//   - executor: REST request executor.
//
// Returns:
//   - Premium-index client.
//
// Version:
//   - 2026-08-19: Added.
func NewClient(executor transport.Executor) (*Client, error) {
	if executor == nil {
		return nil, fmt.Errorf("failed to create USDⓈ-M premium index client: executor=null")
	}
	return &Client{executor: executor}, nil
}

// Send requests mark-price and funding-rate observations.
//
// Notes:
//   - A symbol returns one entry; an empty symbol returns all entries.
//
// Parameters:
//   - ctx: Context for the operation.
//   - params: Request parameters.
//
// Returns:
//   - Premium-index observations normalized as a slice.
//
// Version:
//   - 2026-08-19: Added.
func (c *Client) Send(ctx context.Context, params Params) ([]PremiumIndex, error) {
	if c == nil || c.executor == nil {
		return nil, fmt.Errorf("failed to request USDⓈ-M premium index: client=null")
	}
	query := url.Values{}
	if symbol := strings.TrimSpace(params.Symbol); symbol != "" {
		query.Set("symbol", symbol)
	}
	body, err := c.executor.Do(ctx, transport.Request{Method: http.MethodGet, Path: endpoint.USDSMPremiumIndexPath, Query: query})
	if err != nil {
		return nil, fmt.Errorf("failed to request USDⓈ-M premium index: %w", err)
	}
	result, err := decodeResponse(body)
	if err != nil {
		return nil, fmt.Errorf("failed to request USDⓈ-M premium index: failed to decode response body: %w", err)
	}
	return result, nil
}

func decodeResponse(body []byte) ([]PremiumIndex, error) {
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 {
		return nil, fmt.Errorf("response_body=empty")
	}
	switch trimmed[0] {
	case '[':
		var result []PremiumIndex
		if err := json.Unmarshal(trimmed, &result); err != nil {
			return nil, err
		}
		return result, nil
	case '{':
		var result PremiumIndex
		if err := json.Unmarshal(trimmed, &result); err != nil {
			return nil, err
		}
		return []PremiumIndex{result}, nil
	default:
		return nil, fmt.Errorf("response_body=invalid")
	}
}
