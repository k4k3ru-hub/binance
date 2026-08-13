package exchange_info

import (
	"context"
	"encoding/json"
	"reflect"
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

func TestSymbolUnmarshalJSONNormalizesPermissionSets(t *testing.T) {
	tests := []struct {
		name string
		body string
		want [][]string
	}{
		{
			name: "flat",
			body: `{"symbol":"BTCUSDT","permissionSets":["GRID","COPY"]}`,
			want: [][]string{{"GRID", "COPY"}},
		},
		{
			name: "nested",
			body: `{"symbol":"BTCUSDT","permissionSets":[["GRID"],["COPY","TRADING"]]}`,
			want: [][]string{{"GRID"}, {"COPY", "TRADING"}},
		},
		{
			name: "missing",
			body: `{"symbol":"BTCUSDT"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var symbol Symbol
			if err := json.Unmarshal([]byte(test.body), &symbol); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}
			if !reflect.DeepEqual(symbol.PermissionSets, test.want) {
				t.Fatalf("PermissionSets = %#v, want %#v", symbol.PermissionSets, test.want)
			}
		})
	}
}
