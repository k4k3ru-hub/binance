package markprice

import (
	"context"
	"encoding/json"
	"testing"

	wsprotocol "github.com/k4k3ru-hub/binance/go/websocket/protocol"
)

type recordingExecutor struct {
	key     string
	request wsprotocol.SubscriptionRequest
}

func (e *recordingExecutor) Subscribe(_ context.Context, key string, payload []byte) error {
	e.key = key
	return json.Unmarshal(payload, &e.request)
}
func (e *recordingExecutor) Unsubscribe(_ context.Context, key string, payload []byte) error {
	e.key = key
	return json.Unmarshal(payload, &e.request)
}

func TestParamsStreamSupportsMarkPriceSpeeds(t *testing.T) {
	for _, test := range []struct {
		params Params
		want   string
	}{
		{Params{Symbol: " BTCUSDT "}, "btcusdt@markPrice"},
		{Params{Symbol: "BTCUSDT", UpdateSpeed: UpdateSpeed1s}, "btcusdt@markPrice@1s"},
	} {
		got, err := test.params.Stream()
		if err != nil || got != test.want {
			t.Fatalf("Stream() = %q, %v, want %q", got, err, test.want)
		}
	}
}

func TestParamsStreamRejectsInvalidValues(t *testing.T) {
	for _, params := range []Params{{}, {Symbol: "BTCUSDT", UpdateSpeed: "500ms"}} {
		if _, err := params.Stream(); err == nil {
			t.Fatalf("Stream(%+v) error = nil", params)
		}
	}
}

func TestClientBuildsSubscriptionPayloads(t *testing.T) {
	executor := &recordingExecutor{}
	client, err := NewClient(executor)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	params := Params{Symbol: "ETHUSDT", UpdateSpeed: UpdateSpeed1s}
	if err := client.Subscribe(context.Background(), params); err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}
	if executor.key != "ethusdt@markPrice@1s" || executor.request.Method != wsprotocol.MethodSubscribe || executor.request.ID == 0 {
		t.Fatalf("subscription = %q %#v", executor.key, executor.request)
	}
	if err := client.Unsubscribe(context.Background(), params); err != nil {
		t.Fatalf("Unsubscribe() error = %v", err)
	}
	if executor.request.Method != wsprotocol.MethodUnsubscribe || executor.request.ID == 0 {
		t.Fatalf("unsubscription = %q %#v", executor.key, executor.request)
	}
}
