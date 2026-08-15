// Package bookticker implements Binance USDⓈ-M Futures book-ticker subscriptions.
package bookticker

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	wsprotocol "github.com/k4k3ru-hub/binance/go/websocket/protocol"
	"github.com/k4k3ru-hub/binance/go/websocket/subscriptions"
)

// Params identifies one symbol's best bid/ask stream.
type Params struct {
	Symbol string
}

type Client struct{ executor subscriptions.Executor }

func (p Params) Stream() (string, error) {
	symbol := strings.ToLower(strings.TrimSpace(p.Symbol))
	if symbol == "" {
		return "", fmt.Errorf("failed to build USDⓈ-M book ticker stream: symbol=empty")
	}
	return symbol + "@bookTicker", nil
}

func (p Params) Key() (string, error) { return p.Stream() }

func NewClient(executor subscriptions.Executor) (*Client, error) {
	if executor == nil {
		return nil, fmt.Errorf("failed to create USDⓈ-M book ticker subscription client: executor=null")
	}
	return &Client{executor: executor}, nil
}

func (c *Client) Subscribe(ctx context.Context, params Params) error {
	return c.send(ctx, wsprotocol.MethodSubscribe, params)
}

func (c *Client) Unsubscribe(ctx context.Context, params Params) error {
	return c.send(ctx, wsprotocol.MethodUnsubscribe, params)
}

func (c *Client) send(ctx context.Context, method string, params Params) error {
	operation := strings.ToLower(method)
	if c == nil || c.executor == nil {
		return fmt.Errorf("failed to %s USDⓈ-M book ticker: client=null", operation)
	}
	stream, err := params.Stream()
	if err != nil {
		return fmt.Errorf("failed to %s USDⓈ-M book ticker: %w", operation, err)
	}
	payload, err := json.Marshal(wsprotocol.SubscriptionRequest{
		Method: method, Params: []string{stream}, ID: wsprotocol.RequestID(method, stream),
	})
	if err != nil {
		return fmt.Errorf("failed to %s USDⓈ-M book ticker: failed to encode request: %w", operation, err)
	}
	if method == wsprotocol.MethodSubscribe {
		if err := c.executor.Subscribe(ctx, stream, payload); err != nil {
			return fmt.Errorf("failed to subscribe USDⓈ-M book ticker: %w", err)
		}
		return nil
	}
	if err := c.executor.Unsubscribe(ctx, stream, payload); err != nil {
		return fmt.Errorf("failed to unsubscribe USDⓈ-M book ticker: %w", err)
	}
	return nil
}
