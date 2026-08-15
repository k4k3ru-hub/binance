package bookticker

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

func TestSubscribeAndUnsubscribe(t *testing.T) {
	executor := &recordingExecutor{}
	client, _ := NewClient(executor)
	params := Params{Symbol: "ETHUSDT"}
	if err := client.Subscribe(context.Background(), params); err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}
	if executor.key != "ethusdt@bookTicker" || executor.request.Method != wsprotocol.MethodSubscribe {
		t.Fatalf("subscription = %q %#v", executor.key, executor.request)
	}
	if err := client.Unsubscribe(context.Background(), params); err != nil {
		t.Fatalf("Unsubscribe() error = %v", err)
	}
	if executor.request.Method != wsprotocol.MethodUnsubscribe {
		t.Fatalf("request = %#v", executor.request)
	}
}
