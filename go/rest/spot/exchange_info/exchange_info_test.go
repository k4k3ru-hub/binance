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
	return []byte(`{"timezone":"UTC","serverTime":1,"symbols":[{"symbol":"BTCUSDT","filters":[{"filterType":"PRICE_FILTER","tickSize":"0.01"}]}]}`), nil
}

func TestSendBuildsQueryAndDecodesResponse(t *testing.T) {
	executor := &recordingExecutor{}
	client, err := NewClient(executor)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	show := false
	result, err := client.Send(context.Background(), Params{
		Symbols: []string{"BTCUSDT", "ETHUSDT"}, Permissions: []string{"SPOT"}, ShowPermissionSets: &show,
	})
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if executor.request.Path != endpoint.SpotExchangeInfoPath {
		t.Errorf("path = %q", executor.request.Path)
	}
	if got := executor.request.Query.Get("symbols"); got != `["BTCUSDT","ETHUSDT"]` {
		t.Errorf("symbols = %q", got)
	}
	if got := executor.request.Query.Get("permissions"); got != `["SPOT"]` {
		t.Errorf("permissions = %q", got)
	}
	if got := executor.request.Query.Get("showPermissionSets"); got != "false" {
		t.Errorf("showPermissionSets = %q", got)
	}
	if len(result.Symbols) != 1 || result.Symbols[0].Filters[0].TickSize != "0.01" {
		t.Fatalf("result = %#v", result)
	}
}
