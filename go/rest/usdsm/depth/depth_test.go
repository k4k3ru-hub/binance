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
	return []byte(`{"lastUpdateId":42,"E":100,"T":99,"bids":[["1.0","2.0"]],"asks":[["1.1","3.0"]]}`), nil
}

func TestSendBuildsQueryAndDecodesResponse(t *testing.T) {
	executor := &recordingExecutor{}
	client, _ := NewClient(executor)
	result, err := client.Send(context.Background(), Params{Symbol: "BTCUSDT", Limit: 500})
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if executor.request.Path != endpoint.USDSMDepthPath || executor.request.Query.Get("limit") != "500" {
		t.Fatalf("request = %#v", executor.request)
	}
	if result.MessageTime != 100 || result.TransactionTime != 99 || result.Bids[0].Quantity != "2.0" {
		t.Fatalf("result = %#v", result)
	}
}
