// Package markprice implements Binance USDⓈ-M Futures mark-price subscriptions.
package markprice

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	wsprotocol "github.com/k4k3ru-hub/binance/go/websocket/protocol"
	"github.com/k4k3ru-hub/binance/go/websocket/subscriptions"
)

type UpdateSpeed string

const (
	UpdateSpeedDefault UpdateSpeed = ""
	UpdateSpeed1s      UpdateSpeed = "1s"
)

type Params struct {
	Symbol      string
	UpdateSpeed UpdateSpeed
}

type Client struct{ executor subscriptions.Executor }

// Stream returns the Binance stream name.
//
// Version:
//   - 2026-08-26: Added.
func (p Params) Stream() (string, error) {
	symbol := strings.ToLower(strings.TrimSpace(p.Symbol))
	if symbol == "" {
		return "", fmt.Errorf("failed to build USDⓈ-M mark price stream: symbol=empty")
	}
	stream := symbol + "@markPrice"
	switch p.UpdateSpeed {
	case UpdateSpeedDefault:
	case UpdateSpeed1s:
		stream += "@" + string(p.UpdateSpeed)
	default:
		return "", fmt.Errorf("failed to build USDⓈ-M mark price stream: update_speed=%q", p.UpdateSpeed)
	}
	return stream, nil
}

// Key returns the transport key for the subscription.
//
// Version:
//   - 2026-08-26: Added.
func (p Params) Key() (string, error) { return p.Stream() }

// NewClient creates a USDⓈ-M mark-price subscription client.
//
// Version:
//   - 2026-08-26: Added.
func NewClient(executor subscriptions.Executor) (*Client, error) {
	if executor == nil {
		return nil, fmt.Errorf("failed to create USDⓈ-M mark price subscription client: executor=null")
	}
	return &Client{executor: executor}, nil
}

// Subscribe subscribes to USDⓈ-M mark-price updates.
//
// Version:
//   - 2026-08-26: Added.
func (c *Client) Subscribe(ctx context.Context, params Params) error {
	return c.send(ctx, wsprotocol.MethodSubscribe, params)
}

// Unsubscribe unsubscribes from USDⓈ-M mark-price updates.
//
// Version:
//   - 2026-08-26: Added.
func (c *Client) Unsubscribe(ctx context.Context, params Params) error {
	return c.send(ctx, wsprotocol.MethodUnsubscribe, params)
}

func (c *Client) send(ctx context.Context, method string, params Params) error {
	operation := strings.ToLower(method)
	if c == nil || c.executor == nil {
		return fmt.Errorf("failed to %s USDⓈ-M mark price: client=null", operation)
	}
	stream, err := params.Stream()
	if err != nil {
		return fmt.Errorf("failed to %s USDⓈ-M mark price: %w", operation, err)
	}
	payload, err := json.Marshal(wsprotocol.SubscriptionRequest{Method: method, Params: []string{stream}, ID: wsprotocol.RequestID(method, stream)})
	if err != nil {
		return fmt.Errorf("failed to %s USDⓈ-M mark price: failed to encode request: %w", operation, err)
	}
	if method == wsprotocol.MethodSubscribe {
		if err := c.executor.Subscribe(ctx, stream, payload); err != nil {
			return fmt.Errorf("failed to subscribe USDⓈ-M mark price: %w", err)
		}
		return nil
	}
	if err := c.executor.Unsubscribe(ctx, stream, payload); err != nil {
		return fmt.Errorf("failed to unsubscribe USDⓈ-M mark price: %w", err)
	}
	return nil
}
