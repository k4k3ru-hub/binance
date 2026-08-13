// Package usdsm composes Binance USDⓈ-M Futures REST operations.
package usdsm

import (
	"fmt"

	"github.com/k4k3ru-hub/binance/go/rest/transport"
	"github.com/k4k3ru-hub/binance/go/rest/usdsm/depth"
	"github.com/k4k3ru-hub/binance/go/rest/usdsm/exchange_info"
)

type Client struct {
	exchangeInfo *exchange_info.Client
	depth        *depth.Client
}

func NewClient(executor transport.Executor) (*Client, error) {
	if executor == nil {
		return nil, fmt.Errorf("failed to create USDⓈ-M API client: executor=null")
	}
	exchangeInfoClient, err := exchange_info.NewClient(executor)
	if err != nil {
		return nil, fmt.Errorf("failed to create USDⓈ-M API client: %w", err)
	}
	depthClient, err := depth.NewClient(executor)
	if err != nil {
		return nil, fmt.Errorf("failed to create USDⓈ-M API client: %w", err)
	}
	return &Client{exchangeInfo: exchangeInfoClient, depth: depthClient}, nil
}

func (c *Client) ExchangeInfo() *exchange_info.Client {
	if c == nil {
		return nil
	}
	return c.exchangeInfo
}
func (c *Client) Depth() *depth.Client {
	if c == nil {
		return nil
	}
	return c.depth
}
