// Package websocket provides Binance WebSocket clients.
package websocket

import (
	"context"
	"fmt"
	"net/http"

	spotdepth "github.com/k4k3ru-hub/binance/go/websocket/spot/depth"
	"github.com/k4k3ru-hub/binance/go/websocket/subscriptions"
	usdsmdepth "github.com/k4k3ru-hub/binance/go/websocket/usdsm/depth"
	k4websocket "github.com/k4k3ru-hub/websocket/go"
)

const (
	DefaultSpotEndpointURL  = "wss://stream.binance.com:9443/ws"
	DefaultUSDSMEndpointURL = "wss://fstream.binance.com/ws"
)

type ClientOption = k4websocket.ClientOption
type SessionHandler = k4websocket.SessionHandler
type SessionContext = k4websocket.SessionContext

// SpotClient is a composed Binance Spot WebSocket client.
type SpotClient struct {
	connection *connection
	depth      *spotdepth.Client
}

// USDSMClient is a composed Binance USDⓈ-M Futures WebSocket client.
type USDSMClient struct {
	connection *connection
	depth      *usdsmdepth.Client
}

type connection struct{ client *k4websocket.Client }

// DefaultSpotClientOption returns default Spot WebSocket options.
func DefaultSpotClientOption() *ClientOption {
	option := k4websocket.DefaultClientOption()
	option.EndpointURL = DefaultSpotEndpointURL
	return option
}

// DefaultUSDSMClientOption returns default USDⓈ-M WebSocket options.
func DefaultUSDSMClientOption() *ClientOption {
	option := k4websocket.DefaultClientOption()
	option.EndpointURL = DefaultUSDSMEndpointURL
	return option
}

// NewSpotClient creates and composes a lazy-connecting Spot WebSocket client.
func NewSpotClient(ctx context.Context, handler SessionHandler, option *ClientOption) (*SpotClient, error) {
	conn, err := newConnection(ctx, DefaultSpotEndpointURL, handler, option)
	if err != nil {
		return nil, fmt.Errorf("failed to create Spot WebSocket client: %w", err)
	}
	depthClient, err := spotdepth.NewClient(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to create Spot WebSocket client: %w", err)
	}
	return &SpotClient{connection: conn, depth: depthClient}, nil
}

// NewUSDSMClient creates and composes a lazy-connecting USDⓈ-M WebSocket client.
func NewUSDSMClient(ctx context.Context, handler SessionHandler, option *ClientOption) (*USDSMClient, error) {
	conn, err := newConnection(ctx, DefaultUSDSMEndpointURL, handler, option)
	if err != nil {
		return nil, fmt.Errorf("failed to create USDⓈ-M WebSocket client: %w", err)
	}
	depthClient, err := usdsmdepth.NewClient(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to create USDⓈ-M WebSocket client: %w", err)
	}
	return &USDSMClient{connection: conn, depth: depthClient}, nil
}

func newConnection(ctx context.Context, defaultEndpoint string, handler SessionHandler, option *ClientOption) (*connection, error) {
	if handler == nil {
		return nil, fmt.Errorf("session_handler=null")
	}
	cloned := cloneClientOption(option)
	endpointURL := cloned.EndpointURL
	if endpointURL == "" {
		endpointURL = defaultEndpoint
	}
	client, err := k4websocket.NewClient(ctx, endpointURL, handler, cloned)
	if err != nil {
		return nil, err
	}
	return &connection{client: client}, nil
}

func cloneClientOption(option *ClientOption) *ClientOption {
	if option == nil {
		return k4websocket.DefaultClientOption()
	}
	cloned := &ClientOption{
		EndpointURL: option.EndpointURL, ConnectTimeout: option.ConnectTimeout,
		HandshakeTimeout: option.HandshakeTimeout,
	}
	if option.HTTPHeader != nil {
		cloned.HTTPHeader = option.HTTPHeader.Clone()
	} else {
		cloned.HTTPHeader = make(http.Header)
	}
	if option.SessionOption != nil {
		cloned.SessionOption = option.SessionOption.Clone()
	} else {
		cloned.SessionOption = k4websocket.DefaultSessionOption()
	}
	return cloned
}

func (c *connection) Subscribe(ctx context.Context, key string, payload []byte) error {
	if c == nil || c.client == nil {
		return fmt.Errorf("failed to subscribe WebSocket: connection=null")
	}
	return c.client.Subscribe(ctx, key, payload)
}

func (c *connection) Unsubscribe(ctx context.Context, key string, payload []byte) error {
	if c == nil || c.client == nil {
		return fmt.Errorf("failed to unsubscribe WebSocket: connection=null")
	}
	return c.client.Unsubscribe(ctx, key, payload)
}

var _ subscriptions.Executor = (*connection)(nil)

// Depth returns the Spot diff-depth subscription client.
func (c *SpotClient) Depth() *spotdepth.Client {
	if c == nil {
		return nil
	}
	return c.depth
}

// Connect establishes the Spot WebSocket connection before the first subscription.
func (c *SpotClient) Connect(ctx context.Context) error {
	if c == nil || c.connection == nil {
		return fmt.Errorf("failed to connect Spot WebSocket: client=null")
	}
	return c.connection.client.Connect(ctx)
}

// Close permanently closes the Spot WebSocket client.
func (c *SpotClient) Close() error {
	if c == nil || c.connection == nil {
		return fmt.Errorf("failed to close Spot WebSocket: client=null")
	}
	return c.connection.client.Close()
}

// SessionContext returns the active Spot session, connecting lazily if needed.
func (c *SpotClient) SessionContext() (SessionContext, error) {
	if c == nil || c.connection == nil {
		return nil, fmt.Errorf("failed to get Spot session: client=null")
	}
	return c.connection.client.SessionContext()
}

// Depth returns the USDⓈ-M diff-depth subscription client.
func (c *USDSMClient) Depth() *usdsmdepth.Client {
	if c == nil {
		return nil
	}
	return c.depth
}

// Connect establishes the USDⓈ-M WebSocket connection before the first subscription.
func (c *USDSMClient) Connect(ctx context.Context) error {
	if c == nil || c.connection == nil {
		return fmt.Errorf("failed to connect USDⓈ-M WebSocket: client=null")
	}
	return c.connection.client.Connect(ctx)
}

// Close permanently closes the USDⓈ-M WebSocket client.
func (c *USDSMClient) Close() error {
	if c == nil || c.connection == nil {
		return fmt.Errorf("failed to close USDⓈ-M WebSocket: client=null")
	}
	return c.connection.client.Close()
}

// SessionContext returns the active USDⓈ-M session, connecting lazily if needed.
func (c *USDSMClient) SessionContext() (SessionContext, error) {
	if c == nil || c.connection == nil {
		return nil, fmt.Errorf("failed to get USDⓈ-M session: client=null")
	}
	return c.connection.client.SessionContext()
}
