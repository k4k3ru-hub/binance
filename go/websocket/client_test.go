package websocket

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	gorillawebsocket "github.com/gorilla/websocket"
	wsprotocol "github.com/k4k3ru-hub/binance/go/websocket/protocol"
	spotbookticker "github.com/k4k3ru-hub/binance/go/websocket/spot/book_ticker"
)

type testHandler struct{}

func (*testHandler) HandleMessage(SessionContext, []byte) {}
func (*testHandler) HandleClose(SessionContext)           {}

func TestConstructorsComposeSubscriptionClientsWithoutConnecting(t *testing.T) {
	spot, err := NewSpotClient(context.Background(), &testHandler{}, nil)
	if err != nil || spot == nil || spot.Depth() == nil || spot.BookTicker() == nil || spot.Trades() == nil {
		t.Fatalf("NewSpotClient() = %#v, %v", spot, err)
	}
	futures, err := NewUSDSMClient(context.Background(), &testHandler{}, nil)
	if err != nil || futures == nil || futures.Depth() == nil || futures.BookTicker() == nil || futures.Trades() == nil {
		t.Fatalf("NewUSDSMClient() = %#v, %v", futures, err)
	}
}

func TestConstructorsRequireHandler(t *testing.T) {
	if _, err := NewSpotClient(context.Background(), nil, nil); err == nil {
		t.Fatal("NewSpotClient() error = nil")
	}
	if _, err := NewUSDSMClient(context.Background(), nil, nil); err == nil {
		t.Fatal("NewUSDSMClient() error = nil")
	}
}

type reconnectHandler struct{ closed chan struct{} }

func (*reconnectHandler) HandleMessage(SessionContext, []byte) {}
func (h *reconnectHandler) HandleClose(SessionContext) {
	select {
	case h.closed <- struct{}{}:
	default:
	}
}

func TestSpotBookTickerSubscriptionIsRestoredAfterReconnect(t *testing.T) {
	messages := make(chan wsprotocol.SubscriptionRequest, 2)
	connections := make(chan struct{}, 2)
	upgrader := gorillawebsocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		conn, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		connections <- struct{}{}
		_, payload, err := conn.ReadMessage()
		if err != nil {
			return
		}
		var subscription wsprotocol.SubscriptionRequest
		if json.Unmarshal(payload, &subscription) == nil {
			messages <- subscription
		}
		if len(connections) == 1 {
			_ = conn.WriteControl(gorillawebsocket.CloseMessage, gorillawebsocket.FormatCloseMessage(gorillawebsocket.CloseNormalClosure, "reconnect"), time.Now().Add(time.Second))
			return
		}
		<-request.Context().Done()
	}))
	defer server.Close()

	handler := &reconnectHandler{closed: make(chan struct{}, 1)}
	option := DefaultSpotClientOption()
	option.EndpointURL = "ws" + strings.TrimPrefix(server.URL, "http")
	client, err := NewSpotClient(context.Background(), handler, option)
	if err != nil {
		t.Fatalf("NewSpotClient() error = %v", err)
	}
	defer client.Close()
	if err := client.BookTicker().Subscribe(context.Background(), spotbookticker.Params{Symbol: "BTCUSDT"}); err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}

	select {
	case <-handler.closed:
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for first connection to close")
	}
	if err := client.connection.resubscribeAll(context.Background()); err != nil {
		t.Fatalf("ResubscribeAll() error = %v", err)
	}

	for index := 0; index < 2; index++ {
		select {
		case request := <-messages:
			if request.Method != wsprotocol.MethodSubscribe || len(request.Params) != 1 || request.Params[0] != "btcusdt@bookTicker" {
				t.Fatalf("message[%d] = %#v", index, request)
			}
		case <-time.After(3 * time.Second):
			t.Fatalf("timed out waiting for subscription message %d", index)
		}
	}
}
