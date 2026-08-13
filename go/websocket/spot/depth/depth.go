// Package depth implements Binance Spot diff-depth subscriptions.
package depth

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
	UpdateSpeed100ms   UpdateSpeed = "100ms"
	UpdateSpeed1000ms  UpdateSpeed = "1000ms"
)

type Params struct {
	Symbol      string
	UpdateSpeed UpdateSpeed
}

type Client struct{ executor subscriptions.Executor }

func (p Params) Stream() (string, error) {
	symbol := strings.ToLower(strings.TrimSpace(p.Symbol))
	if symbol == "" {
		return "", fmt.Errorf("failed to build Spot depth stream: symbol=empty")
	}
	stream := symbol + "@depth"
	switch p.UpdateSpeed {
	case UpdateSpeedDefault, UpdateSpeed1000ms:
	case UpdateSpeed100ms:
		stream += "@" + string(p.UpdateSpeed)
	default:
		return "", fmt.Errorf("failed to build Spot depth stream: update_speed=%q", p.UpdateSpeed)
	}
	return stream, nil
}

func (p Params) Key() (string, error) { return p.Stream() }

func NewClient(executor subscriptions.Executor) (*Client, error) {
	if executor == nil {
		return nil, fmt.Errorf("failed to create Spot depth subscription client: executor=null")
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
	if c == nil || c.executor == nil {
		return fmt.Errorf("failed to %s Spot depth: client=null", strings.ToLower(method))
	}
	stream, err := params.Stream()
	if err != nil {
		return fmt.Errorf("failed to %s Spot depth: %w", strings.ToLower(method), err)
	}
	payload, err := json.Marshal(wsprotocol.SubscriptionRequest{Method: method, Params: []string{stream}, ID: wsprotocol.RequestID(method, stream)})
	if err != nil {
		return fmt.Errorf("failed to %s Spot depth: failed to encode request: %w", strings.ToLower(method), err)
	}
	if method == wsprotocol.MethodSubscribe {
		if err := c.executor.Subscribe(ctx, stream, payload); err != nil {
			return fmt.Errorf("failed to subscribe Spot depth: %w", err)
		}
		return nil
	}
	if err := c.executor.Unsubscribe(ctx, stream, payload); err != nil {
		return fmt.Errorf("failed to unsubscribe Spot depth: %w", err)
	}
	return nil
}
