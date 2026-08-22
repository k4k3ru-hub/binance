package websocket

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	spotbookticker "github.com/k4k3ru-hub/binance/go/websocket/spot/book_ticker"
	spottrades "github.com/k4k3ru-hub/binance/go/websocket/spot/trades"
	usdsmbookticker "github.com/k4k3ru-hub/binance/go/websocket/usdsm/book_ticker"
	usdsmtrades "github.com/k4k3ru-hub/binance/go/websocket/usdsm/trades"
)

type fakeSession struct{ id uint64 }

func (s *fakeSession) ID() uint64       { return s.id }
func (*fakeSession) Close()             {}
func (*fakeSession) Send([]byte) error  { return nil }
func (*fakeSession) SendJSON(any) error { return nil }

type fakeWebSocketClient struct {
	mu        sync.Mutex
	session   SessionContext
	messages  [][]byte
	sendError error
	closed    bool
}

func (*fakeWebSocketClient) Connect(context.Context) error { return nil }
func (c *fakeWebSocketClient) SessionContext() (SessionContext, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.session, nil
}
func (c *fakeWebSocketClient) SendRaw(_ context.Context, payload []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.sendError != nil {
		return c.sendError
	}
	c.messages = append(c.messages, append([]byte(nil), payload...))
	return nil
}
func (c *fakeWebSocketClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closed = true
	return nil
}

type recordingLimiter struct {
	mu    sync.Mutex
	waits int
}

func (l *recordingLimiter) Wait(ctx context.Context) error {
	l.mu.Lock()
	l.waits++
	l.mu.Unlock()
	return ctx.Err()
}

type blockingLimiter struct{ started chan struct{} }

func (l *blockingLimiter) Wait(ctx context.Context) error {
	close(l.started)
	<-ctx.Done()
	return ctx.Err()
}

func newTestConnection(client websocketClient, factory func() controlMessageLimiter) *connection {
	closeCtx, cancel := context.WithCancel(context.Background())
	return &connection{
		client:         client,
		subscriptions:  make(map[string][]byte),
		newLimiter:     factory,
		closeCtx:       closeCtx,
		cancelCloseCtx: cancel,
	}
}

func TestConnectionSubscribeAndUnsubscribeShareLimiterAndSkipNoOps(t *testing.T) {
	transport := &fakeWebSocketClient{session: &fakeSession{id: 1}}
	limiter := &recordingLimiter{}
	connection := newTestConnection(transport, func() controlMessageLimiter { return limiter })

	if err := connection.Subscribe(context.Background(), "btc", []byte("subscribe")); err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}
	if err := connection.Subscribe(context.Background(), "btc", []byte("duplicate")); err != nil {
		t.Fatalf("duplicate Subscribe() error = %v", err)
	}
	if err := connection.Unsubscribe(context.Background(), "missing", []byte("missing")); err != nil {
		t.Fatalf("missing Unsubscribe() error = %v", err)
	}
	if err := connection.Unsubscribe(context.Background(), "btc", []byte("unsubscribe")); err != nil {
		t.Fatalf("Unsubscribe() error = %v", err)
	}

	if limiter.waits != 2 {
		t.Fatalf("limiter waits = %d, want 2", limiter.waits)
	}
	if len(transport.messages) != 2 {
		t.Fatalf("messages = %d, want 2", len(transport.messages))
	}
}

func TestConnectionDoesNotStoreSubscriptionAfterSendFailure(t *testing.T) {
	transport := &fakeWebSocketClient{session: &fakeSession{id: 1}, sendError: errors.New("send failed")}
	limiter := &recordingLimiter{}
	connection := newTestConnection(transport, func() controlMessageLimiter { return limiter })
	if err := connection.Subscribe(context.Background(), "btc", []byte("subscribe")); err == nil {
		t.Fatal("Subscribe() error = nil")
	}
	if len(connection.subscriptions) != 0 {
		t.Fatalf("subscriptions = %d, want 0", len(connection.subscriptions))
	}
	transport.sendError = nil
	if err := connection.Subscribe(context.Background(), "btc", []byte("subscribe")); err != nil {
		t.Fatalf("retry Subscribe() error = %v", err)
	}
	if limiter.waits != 2 {
		t.Fatalf("limiter waits = %d, want failed send slot to remain consumed", limiter.waits)
	}
}

func TestConnectionCreatesIndependentLimiterForEachPhysicalSessionAndResubscribe(t *testing.T) {
	transport := &fakeWebSocketClient{session: &fakeSession{id: 1}}
	var limiters []*recordingLimiter
	connection := newTestConnection(transport, func() controlMessageLimiter {
		limiter := &recordingLimiter{}
		limiters = append(limiters, limiter)
		return limiter
	})
	if err := connection.Subscribe(context.Background(), "btc", []byte("btc")); err != nil {
		t.Fatalf("Subscribe(btc) error = %v", err)
	}
	if err := connection.Subscribe(context.Background(), "eth", []byte("eth")); err != nil {
		t.Fatalf("Subscribe(eth) error = %v", err)
	}
	transport.session = &fakeSession{id: 2}
	if err := connection.resubscribeAll(context.Background()); err != nil {
		t.Fatalf("resubscribeAll() error = %v", err)
	}
	if len(limiters) != 2 || limiters[0].waits != 2 || limiters[1].waits != 2 {
		t.Fatalf("limiter waits = %#v, want two independent sessions with two waits", limiters)
	}
}

func TestConnectionsOwnIndependentLimiters(t *testing.T) {
	firstLimiter := &recordingLimiter{}
	secondLimiter := &recordingLimiter{}
	first := newTestConnection(&fakeWebSocketClient{session: &fakeSession{id: 1}}, func() controlMessageLimiter { return firstLimiter })
	second := newTestConnection(&fakeWebSocketClient{session: &fakeSession{id: 2}}, func() controlMessageLimiter { return secondLimiter })
	if err := first.Subscribe(context.Background(), "spot", []byte("spot")); err != nil {
		t.Fatalf("Spot Subscribe() error = %v", err)
	}
	if err := second.Subscribe(context.Background(), "usdsm", []byte("usdsm")); err != nil {
		t.Fatalf("USDSM Subscribe() error = %v", err)
	}
	if firstLimiter.waits != 1 || secondLimiter.waits != 1 {
		t.Fatalf("waits = %d, %d; want independent 1, 1", firstLimiter.waits, secondLimiter.waits)
	}
}

func TestSpotStreamClientsShareOneConnectionLimiter(t *testing.T) {
	transport := &fakeWebSocketClient{session: &fakeSession{id: 1}}
	limiter := &recordingLimiter{}
	connection := newTestConnection(transport, func() controlMessageLimiter { return limiter })
	bookTicker, err := spotbookticker.NewClient(connection)
	if err != nil {
		t.Fatalf("bookticker.NewClient() error = %v", err)
	}
	trades, err := spottrades.NewClient(connection)
	if err != nil {
		t.Fatalf("trades.NewClient() error = %v", err)
	}
	if err := bookTicker.Subscribe(context.Background(), spotbookticker.Params{Symbol: "BTCUSDT"}); err != nil {
		t.Fatalf("BookTicker Subscribe() error = %v", err)
	}
	if err := trades.Subscribe(context.Background(), spottrades.Params{Symbol: "BTCUSDT"}); err != nil {
		t.Fatalf("Trades Subscribe() error = %v", err)
	}
	if limiter.waits != 2 {
		t.Fatalf("limiter waits = %d, want 2", limiter.waits)
	}
}

func TestUSDSMStreamClientsShareOneConnectionLimiter(t *testing.T) {
	transport := &fakeWebSocketClient{session: &fakeSession{id: 1}}
	limiter := &recordingLimiter{}
	connection := newTestConnection(transport, func() controlMessageLimiter { return limiter })
	bookTicker, err := usdsmbookticker.NewClient(connection)
	if err != nil {
		t.Fatalf("bookticker.NewClient() error = %v", err)
	}
	trades, err := usdsmtrades.NewClient(connection)
	if err != nil {
		t.Fatalf("trades.NewClient() error = %v", err)
	}
	if err := bookTicker.Subscribe(context.Background(), usdsmbookticker.Params{Symbol: "BTCUSDT"}); err != nil {
		t.Fatalf("BookTicker Subscribe() error = %v", err)
	}
	if err := trades.Subscribe(context.Background(), usdsmtrades.Params{Symbol: "BTCUSDT"}); err != nil {
		t.Fatalf("Trades Subscribe() error = %v", err)
	}
	if limiter.waits != 2 {
		t.Fatalf("limiter waits = %d, want 2", limiter.waits)
	}
}

func TestClientCloseCancelsLimiterWaitWithoutDeadlock(t *testing.T) {
	transport := &fakeWebSocketClient{session: &fakeSession{id: 1}}
	limiter := &blockingLimiter{started: make(chan struct{})}
	connection := newTestConnection(transport, func() controlMessageLimiter { return limiter })
	client := &SpotClient{connection: connection}
	done := make(chan error, 1)
	go func() { done <- connection.Subscribe(context.Background(), "btc", []byte("btc")) }()
	<-limiter.started
	if err := client.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Subscribe() error = %v, want context.Canceled", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Subscribe() remained blocked after Close()")
	}
}
