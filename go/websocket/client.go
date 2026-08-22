// Package websocket provides Binance WebSocket clients.
package websocket

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	spotbookticker "github.com/k4k3ru-hub/binance/go/websocket/spot/book_ticker"
	spotdepth "github.com/k4k3ru-hub/binance/go/websocket/spot/depth"
	spottrades "github.com/k4k3ru-hub/binance/go/websocket/spot/trades"
	"github.com/k4k3ru-hub/binance/go/websocket/subscriptions"
	usdsmbookticker "github.com/k4k3ru-hub/binance/go/websocket/usdsm/book_ticker"
	usdsmdepth "github.com/k4k3ru-hub/binance/go/websocket/usdsm/depth"
	usdsmtrades "github.com/k4k3ru-hub/binance/go/websocket/usdsm/trades"
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
	bookTicker *spotbookticker.Client
	trades     *spottrades.Client
}

// USDSMClient is a composed Binance USDⓈ-M Futures WebSocket client.
type USDSMClient struct {
	connection *connection
	depth      *usdsmdepth.Client
	bookTicker *usdsmbookticker.Client
	trades     *usdsmtrades.Client
}

type websocketClient interface {
	Connect(context.Context) error
	SessionContext() (SessionContext, error)
	SendRaw(context.Context, []byte) error
	Close() error
}

type connection struct {
	client websocketClient

	controlMu         sync.Mutex
	subscriptions     map[string][]byte
	subscriptionOrder []string
	sessionID         uint64
	limiter           controlMessageLimiter
	newLimiter        func() controlMessageLimiter
	closeCtx          context.Context
	cancelCloseCtx    context.CancelFunc
}

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
	bookTickerClient, err := spotbookticker.NewClient(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to create Spot WebSocket client: %w", err)
	}
	tradesClient, err := spottrades.NewClient(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to create Spot WebSocket client: %w", err)
	}
	return &SpotClient{connection: conn, depth: depthClient, bookTicker: bookTickerClient, trades: tradesClient}, nil
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
	bookTickerClient, err := usdsmbookticker.NewClient(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to create USDⓈ-M WebSocket client: %w", err)
	}
	tradesClient, err := usdsmtrades.NewClient(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to create USDⓈ-M WebSocket client: %w", err)
	}
	return &USDSMClient{connection: conn, depth: depthClient, bookTicker: bookTickerClient, trades: tradesClient}, nil
}

func newConnection(ctx context.Context, defaultEndpoint string, handler SessionHandler, option *ClientOption) (*connection, error) {
	if handler == nil {
		return nil, fmt.Errorf("session_handler=null")
	}
	if ctx == nil {
		ctx = context.Background()
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
	closeCtx, cancelCloseCtx := context.WithCancel(ctx)
	return &connection{
		client:        client,
		subscriptions: make(map[string][]byte),
		newLimiter: func() controlMessageLimiter {
			return newControlMessageLimiter(defaultControlMessageInterval)
		},
		closeCtx:       closeCtx,
		cancelCloseCtx: cancelCloseCtx,
	}, nil
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
	if ctx == nil {
		ctx = context.Background()
	}

	c.controlMu.Lock()
	defer c.controlMu.Unlock()
	if _, exists := c.subscriptions[key]; exists {
		return nil
	}
	if err := c.sendControlMessage(ctx, payload); err != nil {
		return fmt.Errorf("failed to subscribe WebSocket: %w", err)
	}
	c.subscriptions[key] = append([]byte(nil), payload...)
	c.subscriptionOrder = append(c.subscriptionOrder, key)
	return nil
}

func (c *connection) Unsubscribe(ctx context.Context, key string, payload []byte) error {
	if c == nil || c.client == nil {
		return fmt.Errorf("failed to unsubscribe WebSocket: connection=null")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	c.controlMu.Lock()
	defer c.controlMu.Unlock()
	if _, exists := c.subscriptions[key]; !exists {
		return nil
	}
	if err := c.sendControlMessage(ctx, payload); err != nil {
		return fmt.Errorf("failed to unsubscribe WebSocket: %w", err)
	}
	delete(c.subscriptions, key)
	for index, subscriptionKey := range c.subscriptionOrder {
		if subscriptionKey == key {
			c.subscriptionOrder = append(c.subscriptionOrder[:index], c.subscriptionOrder[index+1:]...)
			break
		}
	}
	return nil
}

func (c *connection) sendControlMessage(ctx context.Context, payload []byte) error {
	if err := c.client.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect WebSocket: %w", err)
	}
	session, err := c.client.SessionContext()
	if err != nil {
		return fmt.Errorf("failed to get WebSocket session: %w", err)
	}
	if session == nil {
		return fmt.Errorf("failed to get WebSocket session: session=null")
	}
	if session.ID() != c.sessionID {
		c.sessionID = session.ID()
		c.limiter = c.newLimiter()
	}

	waitCtx, cancel := context.WithCancel(ctx)
	stop := context.AfterFunc(c.closeCtx, cancel)
	defer func() {
		stop()
		cancel()
	}()
	if err := c.limiter.Wait(waitCtx); err != nil {
		return fmt.Errorf("failed to wait for WebSocket control message limit: %w", err)
	}
	if err := waitCtx.Err(); err != nil {
		return fmt.Errorf("failed to wait for WebSocket control message limit: %w", err)
	}
	if err := c.client.SendRaw(ctx, payload); err != nil {
		return fmt.Errorf("failed to enqueue WebSocket control message: %w", err)
	}
	return nil
}

func (c *connection) resubscribeAll(ctx context.Context) error {
	if c == nil || c.client == nil {
		return fmt.Errorf("failed to resubscribe WebSocket: connection=null")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	c.controlMu.Lock()
	defer c.controlMu.Unlock()
	for _, key := range c.subscriptionOrder {
		payload := c.subscriptions[key]
		if err := c.sendControlMessage(ctx, payload); err != nil {
			return fmt.Errorf("failed to resubscribe WebSocket: %w", err)
		}
	}
	return nil
}

var _ subscriptions.Executor = (*connection)(nil)

// Depth returns the Spot diff-depth subscription client.
func (c *SpotClient) Depth() *spotdepth.Client {
	if c == nil {
		return nil
	}
	return c.depth
}

// BookTicker returns the Spot best bid/ask subscription client.
func (c *SpotClient) BookTicker() *spotbookticker.Client {
	if c == nil {
		return nil
	}
	return c.bookTicker
}

// Trades returns the Spot trade subscription client.
//
// Version:
//   - 2026-08-16: Added.
func (c *SpotClient) Trades() *spottrades.Client {
	if c == nil {
		return nil
	}
	return c.trades
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
	c.connection.cancelCloseCtx()
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

// BookTicker returns the USDⓈ-M best bid/ask subscription client.
func (c *USDSMClient) BookTicker() *usdsmbookticker.Client {
	if c == nil {
		return nil
	}
	return c.bookTicker
}

// Trades returns the USDⓈ-M aggregate-trade subscription client.
//
// Version:
//   - 2026-08-16: Added.
func (c *USDSMClient) Trades() *usdsmtrades.Client {
	if c == nil {
		return nil
	}
	return c.trades
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
	c.connection.cancelCloseCtx()
	return c.connection.client.Close()
}

// SessionContext returns the active USDⓈ-M session, connecting lazily if needed.
func (c *USDSMClient) SessionContext() (SessionContext, error) {
	if c == nil || c.connection == nil {
		return nil, fmt.Errorf("failed to get USDⓈ-M session: client=null")
	}
	return c.connection.client.SessionContext()
}
