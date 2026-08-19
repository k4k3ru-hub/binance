package premiumindex

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

const responseObject = `{"symbol":"BTCUSDT","markPrice":"11793.63104562","indexPrice":"11781.80495970","estimatedSettlePrice":"11781.16138815","lastFundingRate":"0.00038246","interestRate":"0.00010000","nextFundingTime":1597392000000,"time":1597370495002}`

func TestSendRequestsSymbolAndNormalizesObjectResponse(t *testing.T) {
	executor := &recordingExecutor{body: []byte(responseObject)}
	client, err := NewClient(executor)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	result, err := client.Send(context.Background(), Params{Symbol: " BTCUSDT "})
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if executor.request.Method != "GET" || executor.request.Path != endpoint.USDSMPremiumIndexPath || executor.request.Query.Get("symbol") != "BTCUSDT" {
		t.Fatalf("request = %#v", executor.request)
	}
	if len(result) != 1 || result[0].LastFundingRate != "0.00038246" || result[0].NextFundingTime != 1597392000000 || result[0].IndexPrice != "11781.80495970" {
		t.Fatalf("result = %#v", result)
	}
}

func TestSendDecodesAllSymbolsResponse(t *testing.T) {
	executor := &recordingExecutor{body: []byte(`[` + responseObject + `]`)}
	client, _ := NewClient(executor)
	result, err := client.Send(context.Background(), Params{})
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if executor.request.Query.Has("symbol") || len(result) != 1 || result[0].Symbol != "BTCUSDT" {
		t.Fatalf("request = %#v, result = %#v", executor.request, result)
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

func TestSendRejectsInvalidResponsesAndNilClient(t *testing.T) {
	for _, body := range []string{"", "{", "null"} {
		client, _ := NewClient(&recordingExecutor{body: []byte(body)})
		if _, err := client.Send(context.Background(), Params{}); err == nil || !strings.Contains(err.Error(), "failed to decode response body") {
			t.Fatalf("Send() body %q error = %v", body, err)
		}
	}
	var nilClient *Client
	if _, err := nilClient.Send(context.Background(), Params{}); err == nil || !strings.Contains(err.Error(), "client=null") {
		t.Fatalf("Send() nil client error = %v", err)
	}
}
