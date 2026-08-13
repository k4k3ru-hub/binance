// Package depth implements the Binance Spot order-book endpoint.
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
	Symbol       string
	Limit        int
	SymbolStatus string
}
type Depth struct {
	LastUpdateID int64                 `json:"lastUpdateId"`
	Bids         []protocol.PriceLevel `json:"bids"`
	Asks         []protocol.PriceLevel `json:"asks"`
}

func NewClient(executor transport.Executor) (*Client, error) {
	if executor == nil {
		return nil, fmt.Errorf("failed to create Spot depth client: executor=null")
	}
	return &Client{executor: executor}, nil
}

func (c *Client) Send(ctx context.Context, params Params) (*Depth, error) {
	if c == nil {
		return nil, fmt.Errorf("failed to request Spot depth: client=null")
	}
	if params.Symbol == "" {
		return nil, fmt.Errorf("failed to request Spot depth: symbol=empty")
	}
	query := url.Values{"symbol": {params.Symbol}}
	if params.Limit != 0 {
		query.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.SymbolStatus != "" {
		query.Set("symbolStatus", params.SymbolStatus)
	}
	body, err := c.executor.Do(ctx, transport.Request{Method: http.MethodGet, Path: endpoint.SpotDepthPath, Query: query})
	if err != nil {
		return nil, fmt.Errorf("failed to request Spot depth: %w", err)
	}
	var result Depth
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to request Spot depth: failed to decode response body: %w", err)
	}
	return &result, nil
}
