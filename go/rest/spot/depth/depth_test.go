package depth

import (
	"context"
	"testing"

	"github.com/k4k3ru-hub/binance/go/rest/endpoint"
	"github.com/k4k3ru-hub/binance/go/rest/transport"
)

type recordingExecutor struct{ request transport.Request }

func (e *recordingExecutor) Do(_ context.Context, request transport.Request) ([]byte, error) {
	e.request = request
	return []byte(`{"lastUpdateId":42,"bids":[["1.0","2.0"]],"asks":[["1.1","3.0"]]}`), nil
}

func TestSendBuildsQueryAndDecodesLevels(t *testing.T) {
	executor := &recordingExecutor{}
	client, _ := NewClient(executor)
	result, err := client.Send(context.Background(), Params{Symbol: "BTCUSDT", Limit: 100, SymbolStatus: "TRADING"})
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if executor.request.Path != endpoint.SpotDepthPath || executor.request.Query.Get("symbol") != "BTCUSDT" || executor.request.Query.Get("limit") != "100" {
		t.Fatalf("request = %#v", executor.request)
	}
	if result.LastUpdateID != 42 || result.Bids[0].Price != "1.0" || result.Asks[0].Quantity != "3.0" {
		t.Fatalf("result = %#v", result)
	}
}

func TestSendRejectsEmptySymbol(t *testing.T) {
	client, _ := NewClient(&recordingExecutor{})
	if _, err := client.Send(context.Background(), Params{}); err == nil {
		t.Fatal("Send() error = nil")
	}
}
