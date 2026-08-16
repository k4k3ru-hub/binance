// Package trades implements Binance USDⓈ-M Futures aggregate-trade subscriptions.
package trades

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	wsprotocol "github.com/k4k3ru-hub/binance/go/websocket/protocol"
	"github.com/k4k3ru-hub/binance/go/websocket/subscriptions"
)

// Params identifies one symbol's aggregate-trade stream.
type Params struct {
	Symbol string
}

// Client manages USDⓈ-M aggregate-trade subscriptions.
type Client struct{ executor subscriptions.Executor }

// Stream returns the Binance stream name.
//
// Version:
//   - 2026-08-16: Added.
func (p Params) Stream() (string, error) {
	symbol := strings.ToLower(strings.TrimSpace(p.Symbol))
	if symbol == "" {
		return "", fmt.Errorf("failed to build USDⓈ-M trades stream: symbol=empty")
	}
	return symbol + "@aggTrade", nil
}

// Key returns the transport key for the subscription.
//
// Version:
//   - 2026-08-16: Added.
func (p Params) Key() (string, error) { return p.Stream() }

// NewClient creates a USDⓈ-M aggregate-trade subscription client.
//
// Version:
//   - 2026-08-16: Added.
func NewClient(executor subscriptions.Executor) (*Client, error) {
	if executor == nil {
		return nil, fmt.Errorf("failed to create USDⓈ-M trades subscription client: executor=null")
	}
	return &Client{executor: executor}, nil
}

// Subscribe subscribes to USDⓈ-M aggregate-trade updates.
//
// Version:
//   - 2026-08-16: Added.
func (c *Client) Subscribe(ctx context.Context, params Params) error {
	return c.send(ctx, wsprotocol.MethodSubscribe, params)
}

// Unsubscribe unsubscribes from USDⓈ-M aggregate-trade updates.
//
// Version:
//   - 2026-08-16: Added.
func (c *Client) Unsubscribe(ctx context.Context, params Params) error {
	return c.send(ctx, wsprotocol.MethodUnsubscribe, params)
}

func (c *Client) send(ctx context.Context, method string, params Params) error {
	operation := strings.ToLower(method)
	if c == nil || c.executor == nil {
		return fmt.Errorf("failed to %s USDⓈ-M trades: client=null", operation)
	}
	stream, err := params.Stream()
	if err != nil {
		return fmt.Errorf("failed to %s USDⓈ-M trades: %w", operation, err)
	}
	payload, err := json.Marshal(wsprotocol.SubscriptionRequest{
		Method: method, Params: []string{stream}, ID: wsprotocol.RequestID(method, stream),
	})
	if err != nil {
		return fmt.Errorf("failed to %s USDⓈ-M trades: failed to encode request: %w", operation, err)
	}
	if method == wsprotocol.MethodSubscribe {
		if err := c.executor.Subscribe(ctx, stream, payload); err != nil {
			return fmt.Errorf("failed to subscribe USDⓈ-M trades: %w", err)
		}
		return nil
	}
	if err := c.executor.Unsubscribe(ctx, stream, payload); err != nil {
		return fmt.Errorf("failed to unsubscribe USDⓈ-M trades: %w", err)
	}
	return nil
}
