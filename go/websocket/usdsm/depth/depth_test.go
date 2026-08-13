package depth

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

func TestParamsStreamSupportsFuturesSpeeds(t *testing.T) {
	for speed, want := range map[UpdateSpeed]string{
		UpdateSpeedDefault: "btcusdt@depth", UpdateSpeed100ms: "btcusdt@depth@100ms", UpdateSpeed250ms: "btcusdt@depth", UpdateSpeed500ms: "btcusdt@depth@500ms",
	} {
		got, err := (Params{Symbol: "BTCUSDT", UpdateSpeed: speed}).Stream()
		if err != nil || got != want {
			t.Errorf("Stream(%q) = %q, %v", speed, got, err)
		}
	}
}

func TestClientBuildsSubscriptionPayload(t *testing.T) {
	executor := &recordingExecutor{}
	client, _ := NewClient(executor)
	params := Params{Symbol: "ETHUSDT", UpdateSpeed: UpdateSpeed500ms}
	if err := client.Subscribe(context.Background(), params); err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}
	if executor.key != "ethusdt@depth@500ms" || executor.request.Method != wsprotocol.MethodSubscribe || executor.request.ID == 0 {
		t.Fatalf("subscription = %q %#v", executor.key, executor.request)
	}
}
