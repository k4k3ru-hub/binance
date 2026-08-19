// Package openinteresthistory implements the Binance USDⓈ-M Futures open-interest history endpoint.
package openinteresthistory

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

type Period string

const (
	Period5m  Period = "5m"
	Period15m Period = "15m"
	Period30m Period = "30m"
	Period1h  Period = "1h"
	Period2h  Period = "2h"
	Period4h  Period = "4h"
	Period6h  Period = "6h"
	Period12h Period = "12h"
	Period1d  Period = "1d"
)

type Client struct{ executor transport.Executor }

type Params struct {
	Symbol    string
	Period    Period
	Limit     int
	StartTime int64
	EndTime   int64
}

type OpenInterest struct {
	Symbol               string `json:"symbol"`
	SumOpenInterest      string `json:"sumOpenInterest"`
	SumOpenInterestValue string `json:"sumOpenInterestValue"`
	CMCCirculatingSupply string `json:"CMCCirculatingSupply"`
	Timestamp            int64  `json:"timestamp"`
}

// NewClient creates a USDⓈ-M Futures open-interest history client.
//
// Parameters:
//   - executor: REST request executor.
//
// Returns:
//   - Open-interest history client.
//
// Version:
//   - 2026-08-19: Added.
func NewClient(executor transport.Executor) (*Client, error) {
	if executor == nil {
		return nil, fmt.Errorf("failed to create USDⓈ-M open interest history client: executor=null")
	}
	return &Client{executor: executor}, nil
}

// Send requests historical open-interest observations for a USDⓈ-M Futures symbol.
//
// Parameters:
//   - ctx: Context for the operation.
//   - params: Request parameters.
//
// Returns:
//   - Historical open-interest observations.
//
// Version:
//   - 2026-08-19: Added.
func (c *Client) Send(ctx context.Context, params Params) ([]OpenInterest, error) {
	if c == nil || c.executor == nil {
		return nil, fmt.Errorf("failed to request USDⓈ-M open interest history: client=null")
	}
	query, err := params.query()
	if err != nil {
		return nil, fmt.Errorf("failed to request USDⓈ-M open interest history: %w", err)
	}
	body, err := c.executor.Do(ctx, transport.Request{Method: http.MethodGet, Path: endpoint.USDSMOpenInterestHistoryPath, Query: query})
	if err != nil {
		return nil, fmt.Errorf("failed to request USDⓈ-M open interest history: %w", err)
	}
	var result []OpenInterest
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to request USDⓈ-M open interest history: failed to decode response body: %w", err)
	}
	return result, nil
}

func (p Params) query() (url.Values, error) {
	symbol := strings.TrimSpace(p.Symbol)
	if symbol == "" {
		return nil, fmt.Errorf("failed to validate open interest history parameters: symbol=empty")
	}
	if !p.Period.valid() {
		return nil, fmt.Errorf("failed to validate open interest history parameters: period=invalid")
	}
	if p.Limit < 0 || p.Limit > 500 {
		return nil, fmt.Errorf("failed to validate open interest history parameters: limit=out_of_range min_value=1 max_value=500")
	}
	if p.StartTime < 0 {
		return nil, fmt.Errorf("failed to validate open interest history parameters: start_time=out_of_range min_value=0")
	}
	if p.EndTime < 0 {
		return nil, fmt.Errorf("failed to validate open interest history parameters: end_time=out_of_range min_value=0")
	}
	if p.StartTime != 0 && p.EndTime != 0 && p.StartTime > p.EndTime {
		return nil, fmt.Errorf("failed to validate open interest history parameters: time_range=invalid")
	}

	query := url.Values{"symbol": {symbol}, "period": {string(p.Period)}}
	if p.Limit != 0 {
		query.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.StartTime != 0 {
		query.Set("startTime", strconv.FormatInt(p.StartTime, 10))
	}
	if p.EndTime != 0 {
		query.Set("endTime", strconv.FormatInt(p.EndTime, 10))
	}
	return query, nil
}

func (p Period) valid() bool {
	switch p {
	case Period5m, Period15m, Period30m, Period1h, Period2h, Period4h, Period6h, Period12h, Period1d:
		return true
	default:
		return false
	}
}
