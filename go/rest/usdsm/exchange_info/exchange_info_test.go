package exchange_info

import (
	"context"
	"testing"

	"github.com/k4k3ru-hub/binance/go/rest/endpoint"
	"github.com/k4k3ru-hub/binance/go/rest/transport"
)

type recordingExecutor struct{ request transport.Request }

func (e *recordingExecutor) Do(_ context.Context, request transport.Request) ([]byte, error) {
	e.request = request
	return []byte(`{"timezone":"UTC","serverTime":1,"assets":[{"asset":"USDT","marginAvailable":true}],"symbols":[{"symbol":"BTCUSDT","contractType":"PERPETUAL","filters":[{"filterType":"PERCENT_PRICE","multiplierUp":"1.05","multiplierDecimal":"4"}]}]}`), nil
}

func TestSendRequestsEndpointAndDecodesResponse(t *testing.T) {
	executor := &recordingExecutor{}
	client, _ := NewClient(executor)
	result, err := client.Send(context.Background(), Params{})
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if executor.request.Path != endpoint.USDSMExchangeInfoPath {
		t.Errorf("path = %q", executor.request.Path)
	}
	if len(result.Assets) != 1 || !result.Assets[0].MarginAvailable || result.Symbols[0].Filters[0].MultiplierDecimal != "4" {
		t.Fatalf("result = %#v", result)
	}
}
