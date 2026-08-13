// Package depth implements the Binance USDⓈ-M Futures order-book endpoint.
package depth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/k4k3ru-hub/binance/go/rest/endpoint"
	"github.com/k4k3ru-hub/binance/go/rest/protocol"
	"github.com/k4k3ru-hub/binance/go/rest/transport"
)

type Client struct{ executor transport.Executor }
type Params struct {
	Symbol string
	Limit  int
}
type Depth struct {
	LastUpdateID    int64                 `json:"lastUpdateId"`
	MessageTime     int64                 `json:"E"`
	TransactionTime int64                 `json:"T"`
	Bids            []protocol.PriceLevel `json:"bids"`
	Asks            []protocol.PriceLevel `json:"asks"`
}

func NewClient(executor transport.Executor) (*Client, error) {
	if executor == nil {
		return nil, fmt.Errorf("failed to create USDⓈ-M depth client: executor=null")
	}
	return &Client{executor: executor}, nil
}

func (c *Client) Send(ctx context.Context, params Params) (*Depth, error) {
	if c == nil {
		return nil, fmt.Errorf("failed to request USDⓈ-M depth: client=null")
	}
	if params.Symbol == "" {
		return nil, fmt.Errorf("failed to request USDⓈ-M depth: symbol=empty")
	}
	query := url.Values{"symbol": {params.Symbol}}
	if params.Limit != 0 {
		query.Set("limit", strconv.Itoa(params.Limit))
	}
	body, err := c.executor.Do(ctx, transport.Request{Method: http.MethodGet, Path: endpoint.USDSMDepthPath, Query: query})
	if err != nil {
		return nil, fmt.Errorf("failed to request USDⓈ-M depth: %w", err)
	}
	var result Depth
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to request USDⓈ-M depth: failed to decode response body: %w", err)
	}
	return &result, nil
}
