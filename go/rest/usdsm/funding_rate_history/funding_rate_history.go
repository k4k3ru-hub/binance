// Package fundingratehistory implements the Binance USDⓈ-M Futures funding-rate history endpoint.
package fundingratehistory

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/k4k3ru-hub/binance/go/rest/endpoint"
	"github.com/k4k3ru-hub/binance/go/rest/transport"
)

type RateType string

const (
	RateTypeRegular RateType = "Regular"
	RateTypeSpecial RateType = "Special"
)

type Client struct{ executor transport.Executor }

type Params struct {
	Symbol    string
	StartTime int64
	EndTime   int64
	Limit     int
}

type FundingRate struct {
	Symbol      string   `json:"symbol"`
	FundingRate string   `json:"fundingRate"`
	FundingTime int64    `json:"fundingTime"`
	MarkPrice   string   `json:"markPrice"`
	RateType    RateType `json:"rateType"`
}

// NewClient creates a USDⓈ-M Futures funding-rate history client.
//
// Parameters:
//   - executor: REST request executor.
//
// Returns:
//   - Funding-rate history client.
//
// Version:
//   - 2026-08-19: Added.
func NewClient(executor transport.Executor) (*Client, error) {
	if executor == nil {
		return nil, fmt.Errorf("failed to create USDⓈ-M funding rate history client: executor=null")
	}
	return &Client{executor: executor}, nil
}

// Send requests settled funding-rate history.
//
// Parameters:
//   - ctx: Context for the operation.
//   - params: Request parameters.
//
// Returns:
//   - Funding-rate history in ascending time order.
//
// Version:
//   - 2026-08-19: Added.
func (c *Client) Send(ctx context.Context, params Params) ([]FundingRate, error) {
	if c == nil || c.executor == nil {
		return nil, fmt.Errorf("failed to request USDⓈ-M funding rate history: client=null")
	}
	query, err := params.query()
	if err != nil {
		return nil, fmt.Errorf("failed to request USDⓈ-M funding rate history: %w", err)
	}
	body, err := c.executor.Do(ctx, transport.Request{Method: http.MethodGet, Path: endpoint.USDSMFundingRateHistoryPath, Query: query})
	if err != nil {
		return nil, fmt.Errorf("failed to request USDⓈ-M funding rate history: %w", err)
	}
	var result []FundingRate
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to request USDⓈ-M funding rate history: failed to decode response body: %w", err)
	}
	return result, nil
}

func (p Params) query() (url.Values, error) {
	if p.Limit < 0 || p.Limit > 1000 {
		return nil, fmt.Errorf("failed to validate funding rate history parameters: limit=out_of_range min_value=1 max_value=1000")
	}
	if p.StartTime < 0 {
		return nil, fmt.Errorf("failed to validate funding rate history parameters: start_time=out_of_range min_value=0")
	}
	if p.EndTime < 0 {
		return nil, fmt.Errorf("failed to validate funding rate history parameters: end_time=out_of_range min_value=0")
	}
	if p.StartTime != 0 && p.EndTime != 0 && p.StartTime > p.EndTime {
		return nil, fmt.Errorf("failed to validate funding rate history parameters: time_range=invalid")
	}

	query := url.Values{}
	if symbol := strings.TrimSpace(p.Symbol); symbol != "" {
		query.Set("symbol", symbol)
	}
	if p.StartTime != 0 {
		query.Set("startTime", strconv.FormatInt(p.StartTime, 10))
	}
	if p.EndTime != 0 {
		query.Set("endTime", strconv.FormatInt(p.EndTime, 10))
	}
	if p.Limit != 0 {
		query.Set("limit", strconv.Itoa(p.Limit))
	}
	return query, nil
}
