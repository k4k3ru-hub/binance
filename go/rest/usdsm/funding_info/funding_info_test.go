package fundinginfo

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

func TestSendRequestsEndpointAndDecodesResponse(t *testing.T) {
	executor := &recordingExecutor{body: []byte(`[{"symbol":"BLZUSDT","adjustedFundingRateCap":"0.02500000","adjustedFundingRateFloor":"-0.02500000","fundingIntervalHours":8,"disclaimer":false}]`)}
	client, err := NewClient(executor)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	result, err := client.Send(context.Background(), Params{})
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if executor.request.Method != "GET" || executor.request.Path != endpoint.USDSMFundingInfoPath || len(executor.request.Query) != 0 {
		t.Fatalf("request = %#v", executor.request)
	}
	if len(result) != 1 || result[0].Symbol != "BLZUSDT" || result[0].FundingIntervalHours != 8 || result[0].AdjustedFundingRateFloor != "-0.02500000" {
		t.Fatalf("result = %#v", result)
	}
}

func TestSendWrapsExecutorError(t *testing.T) {
	want := errors.New("transport failure")
	client, _ := NewClient(&recordingExecutor{err: want})
	_, err := client.Send(context.Background(), Params{})
	if !errors.Is(err, want) {
		t.Fatalf("Send() error = %v, want wrapped error", err)
	}
}

func TestSendRejectsInvalidResponseAndNilClient(t *testing.T) {
	client, _ := NewClient(&recordingExecutor{body: []byte(`{`)})
	if _, err := client.Send(context.Background(), Params{}); err == nil || !strings.Contains(err.Error(), "failed to decode response body") {
		t.Fatalf("Send() decode error = %v", err)
	}
	var nilClient *Client
	if _, err := nilClient.Send(context.Background(), Params{}); err == nil || !strings.Contains(err.Error(), "client=null") {
		t.Fatalf("Send() nil client error = %v", err)
	}
}
