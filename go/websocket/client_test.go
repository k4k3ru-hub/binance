package websocket

import (
	"context"
	"testing"
)

type testHandler struct{}

func (*testHandler) HandleMessage(SessionContext, []byte) {}
func (*testHandler) HandleClose(SessionContext)           {}

func TestConstructorsComposeDepthClientsWithoutConnecting(t *testing.T) {
	spot, err := NewSpotClient(context.Background(), &testHandler{}, nil)
	if err != nil || spot == nil || spot.Depth() == nil {
		t.Fatalf("NewSpotClient() = %#v, %v", spot, err)
	}
	futures, err := NewUSDSMClient(context.Background(), &testHandler{}, nil)
	if err != nil || futures == nil || futures.Depth() == nil {
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
