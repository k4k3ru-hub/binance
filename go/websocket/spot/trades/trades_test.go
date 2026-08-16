package trades

import (
	"context"
	"encoding/json"
	"testing"

	wsprotocol "github.com/k4k3ru-hub/binance/go/websocket/protocol"
)

type recordingExecutor struct {
	keys     []string
	requests []wsprotocol.SubscriptionRequest
}

func (e *recordingExecutor) record(key string, payload []byte) error {
	var request wsprotocol.SubscriptionRequest
	if err := json.Unmarshal(payload, &request); err != nil {
		return err
	}
	e.keys = append(e.keys, key)
	e.requests = append(e.requests, request)
	return nil
}

func (e *recordingExecutor) Subscribe(_ context.Context, key string, payload []byte) error {
	return e.record(key, payload)
}

func (e *recordingExecutor) Unsubscribe(_ context.Context, key string, payload []byte) error {
	return e.record(key, payload)
}

func TestSubscribeAndUnsubscribe(t *testing.T) {
	executor := &recordingExecutor{}
	client, err := NewClient(executor)
	if err != nil {
		t.Fatal(err)
	}
	params := Params{Symbol: " BTCUSDT "}
	if err := client.Subscribe(context.Background(), params); err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}
	if err := client.Unsubscribe(context.Background(), params); err != nil {
		t.Fatalf("Unsubscribe() error = %v", err)
	}
	if len(executor.keys) != 2 || executor.keys[0] != "btcusdt@trade" {
		t.Fatalf("keys = %#v", executor.keys)
	}
	if executor.requests[0].Method != wsprotocol.MethodSubscribe || executor.requests[1].Method != wsprotocol.MethodUnsubscribe {
		t.Fatalf("requests = %#v", executor.requests)
	}
	if executor.requests[0].ID == 0 || executor.requests[0].ID == executor.requests[1].ID {
		t.Fatalf("request IDs = %d, %d", executor.requests[0].ID, executor.requests[1].ID)
	}
}

func TestParamsRejectsEmptySymbol(t *testing.T) {
	if _, err := (Params{}).Stream(); err == nil {
		t.Fatal("Stream() error = nil")
	}
}
