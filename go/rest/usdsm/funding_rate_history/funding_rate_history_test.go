package fundingratehistory

import (
	"context"
	"errors"
	"reflect"
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
	executor := &recordingExecutor{body: []byte(`[{"symbol":"BTCUSDT","fundingRate":"-0.03750000","fundingTime":1570608000000,"markPrice":"34287.54619963","rateType":"Regular"}]`)}
	client, err := NewClient(executor)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	result, err := client.Send(context.Background(), Params{Symbol: " BTCUSDT ", StartTime: 1000, EndTime: 2000, Limit: 1000})
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if executor.request.Method != "GET" || executor.request.Path != endpoint.USDSMFundingRateHistoryPath {
		t.Fatalf("request = %#v", executor.request)
	}
	wantQuery := map[string]string{"symbol": "BTCUSDT", "startTime": "1000", "endTime": "2000", "limit": "1000"}
	for key, want := range wantQuery {
		if got := executor.request.Query.Get(key); got != want {
			t.Errorf("query[%q] = %q, want %q", key, got, want)
		}
	}
	want := []FundingRate{{Symbol: "BTCUSDT", FundingRate: "-0.03750000", FundingTime: 1570608000000, MarkPrice: "34287.54619963", RateType: RateTypeRegular}}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("result = %#v, want %#v", result, want)
	}
}

func TestSendOmitsOptionalParameters(t *testing.T) {
	executor := &recordingExecutor{body: []byte(`[]`)}
	client, _ := NewClient(executor)
	if _, err := client.Send(context.Background(), Params{}); err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if len(executor.request.Query) != 0 {
		t.Fatalf("query = %v", executor.request.Query)
	}
}

func TestSendValidatesParameters(t *testing.T) {
	tests := []struct {
		name   string
		params Params
		want   string
	}{
		{name: "negative limit", params: Params{Limit: -1}, want: "limit=out_of_range"},
		{name: "large limit", params: Params{Limit: 1001}, want: "limit=out_of_range"},
		{name: "negative start", params: Params{StartTime: -1}, want: "start_time=out_of_range"},
		{name: "negative end", params: Params{EndTime: -1}, want: "end_time=out_of_range"},
		{name: "reversed range", params: Params{StartTime: 2, EndTime: 1}, want: "time_range=invalid"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, _ := NewClient(&recordingExecutor{})
			_, err := client.Send(context.Background(), test.params)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("Send() error = %v, want containing %q", err, test.want)
			}
		})
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
