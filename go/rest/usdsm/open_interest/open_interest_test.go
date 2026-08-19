package openinterest

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/k4k3ru-hub/binance/go/rest/endpoint"
	"github.com/k4k3ru-hub/binance/go/rest/transport"
)

type recordingExecutor struct {
	request transport.Request
	body    []byte
	err     error
}

func (e *recordingExecutor) Do(_ context.Context, request transport.Request) ([]byte, error) {
	e.request = request
	return e.body, e.err
}

func TestSendBuildsQueryAndDecodesResponse(t *testing.T) {
	executor := &recordingExecutor{body: []byte(`{"openInterest":"10659.509","symbol":"BTCUSDT","time":1589437530011}`)}
	client, err := NewClient(executor)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	result, err := client.Send(context.Background(), Params{Symbol: " BTCUSDT "})
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if executor.request.Method != "GET" || executor.request.Path != endpoint.USDSMOpenInterestPath || executor.request.Query.Get("symbol") != "BTCUSDT" {
		t.Fatalf("request = %#v", executor.request)
	}
	if result.OpenInterest != "10659.509" || result.Symbol != "BTCUSDT" || result.Time != 1589437530011 {
		t.Fatalf("result = %#v", result)
	}
}

func TestSendRejectsEmptySymbol(t *testing.T) {
	client, _ := NewClient(&recordingExecutor{})
	_, err := client.Send(context.Background(), Params{Symbol: " "})
	if err == nil || !strings.Contains(err.Error(), "symbol=empty") {
		t.Fatalf("Send() error = %v", err)
	}
}

func TestSendWrapsExecutorError(t *testing.T) {
	want := errors.New("transport failure")
	client, _ := NewClient(&recordingExecutor{err: want})
	_, err := client.Send(context.Background(), Params{Symbol: "BTCUSDT"})
	if !errors.Is(err, want) {
		t.Fatalf("Send() error = %v, want wrapped error", err)
	}
}

func TestSendRejectsInvalidResponse(t *testing.T) {
	client, _ := NewClient(&recordingExecutor{body: []byte(`{`)})
	_, err := client.Send(context.Background(), Params{Symbol: "BTCUSDT"})
	if err == nil || !strings.Contains(err.Error(), "failed to decode response body") {
		t.Fatalf("Send() error = %v", err)
	}
}

func TestSendRejectsNilClient(t *testing.T) {
	var client *Client
	_, err := client.Send(context.Background(), Params{Symbol: "BTCUSDT"})
	if err == nil || !strings.Contains(err.Error(), "client=null") {
		t.Fatalf("Send() error = %v", err)
	}
}
