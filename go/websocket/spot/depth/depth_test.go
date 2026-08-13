package depth

import (
	"context"
	"encoding/json"
	"sync"
	"testing"

	wsprotocol "github.com/k4k3ru-hub/binance/go/websocket/protocol"
)

type recordingExecutor struct {
	mu           sync.Mutex
	subscribes   map[string]wsprotocol.SubscriptionRequest
	unsubscribes map[string]wsprotocol.SubscriptionRequest
}

func (e *recordingExecutor) Subscribe(_ context.Context, key string, payload []byte) error {
	var request wsprotocol.SubscriptionRequest
	if err := json.Unmarshal(payload, &request); err != nil {
		return err
	}
	e.mu.Lock()
	e.subscribes[key] = request
	e.mu.Unlock()
	return nil
}

func (e *recordingExecutor) Unsubscribe(_ context.Context, key string, payload []byte) error {
	var request wsprotocol.SubscriptionRequest
	if err := json.Unmarshal(payload, &request); err != nil {
		return err
	}
	e.mu.Lock()
	e.unsubscribes[key] = request
	e.mu.Unlock()
	return nil
}

func TestParamsStream(t *testing.T) {
	stream, err := (Params{Symbol: " BTCUSDT ", UpdateSpeed: UpdateSpeed100ms}).Stream()
	if err != nil {
		t.Fatalf("Stream() error = %v", err)
	}
	if stream != "btcusdt@depth@100ms" {
		t.Fatalf("Stream() = %q", stream)
	}
	stream, err = (Params{Symbol: "BTCUSDT", UpdateSpeed: UpdateSpeed1000ms}).Stream()
	if err != nil || stream != "btcusdt@depth" {
		t.Fatalf("Stream(1000ms) = %q, %v", stream, err)
	}
	if _, err := (Params{Symbol: "BTCUSDT", UpdateSpeed: "250ms"}).Stream(); err == nil {
		t.Fatal("unsupported speed error = nil")
	}
}

func TestClientKeepsConcurrentSubscriptionsIndependent(t *testing.T) {
	executor := &recordingExecutor{subscribes: make(map[string]wsprotocol.SubscriptionRequest), unsubscribes: make(map[string]wsprotocol.SubscriptionRequest)}
	client, _ := NewClient(executor)
	var waitGroup sync.WaitGroup
	for _, symbol := range []string{"BTCUSDT", "ETHUSDT"} {
		symbol := symbol
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			params := Params{Symbol: symbol, UpdateSpeed: UpdateSpeed100ms}
			if err := client.Subscribe(context.Background(), params); err != nil {
				t.Errorf("Subscribe() error = %v", err)
			}
			if err := client.Unsubscribe(context.Background(), params); err != nil {
				t.Errorf("Unsubscribe() error = %v", err)
			}
		}()
	}
	waitGroup.Wait()
	for _, key := range []string{"btcusdt@depth@100ms", "ethusdt@depth@100ms"} {
		if executor.subscribes[key].Method != wsprotocol.MethodSubscribe {
			t.Errorf("subscribe[%q] = %#v", key, executor.subscribes[key])
		}
		if executor.unsubscribes[key].Method != wsprotocol.MethodUnsubscribe {
			t.Errorf("unsubscribe[%q] = %#v", key, executor.unsubscribes[key])
		}
	}
}
