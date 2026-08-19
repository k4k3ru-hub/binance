// Package fundinginfo implements the Binance USDⓈ-M Futures funding-rate information endpoint.
package fundinginfo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/k4k3ru-hub/binance/go/rest/endpoint"
	"github.com/k4k3ru-hub/binance/go/rest/transport"
)

type Client struct{ executor transport.Executor }

type Params struct{}

type FundingInfo struct {
	Symbol                   string `json:"symbol"`
	AdjustedFundingRateCap   string `json:"adjustedFundingRateCap"`
	AdjustedFundingRateFloor string `json:"adjustedFundingRateFloor"`
	FundingIntervalHours     int64  `json:"fundingIntervalHours"`
	Disclaimer               bool   `json:"disclaimer"`
}

// NewClient creates a USDⓈ-M Futures funding-information client.
//
// Parameters:
//   - executor: REST request executor.
//
// Returns:
//   - Funding-information client.
//
// Version:
//   - 2026-08-19: Added.
func NewClient(executor transport.Executor) (*Client, error) {
	if executor == nil {
		return nil, fmt.Errorf("failed to create USDⓈ-M funding info client: executor=null")
	}
	return &Client{executor: executor}, nil
}

// Send requests funding information for symbols with adjusted settings.
//
// Parameters:
//   - ctx: Context for the operation.
//   - params: Request parameters.
//
// Returns:
//   - Funding information for adjusted symbols.
//
// Version:
//   - 2026-08-19: Added.
func (c *Client) Send(ctx context.Context, _ Params) ([]FundingInfo, error) {
	if c == nil || c.executor == nil {
		return nil, fmt.Errorf("failed to request USDⓈ-M funding info: client=null")
	}
	body, err := c.executor.Do(ctx, transport.Request{Method: http.MethodGet, Path: endpoint.USDSMFundingInfoPath})
	if err != nil {
		return nil, fmt.Errorf("failed to request USDⓈ-M funding info: %w", err)
	}
	var result []FundingInfo
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to request USDⓈ-M funding info: failed to decode response body: %w", err)
	}
	return result, nil
}
