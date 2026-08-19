package openinteresthistory

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
	executor := &recordingExecutor{body: []byte(`[{"symbol":"BTCUSDT","sumOpenInterest":"20403.12345678","sumOpenInterestValue":"176196512.12345678","CMCCirculatingSupply":"165880.538","timestamp":1591261042378}]`)}
	client, err := NewClient(executor)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	result, err := client.Send(context.Background(), Params{
		Symbol: " BTCUSDT ", Period: Period5m, Limit: 500, StartTime: 1000, EndTime: 2000,
	})
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if executor.request.Method != "GET" || executor.request.Path != endpoint.USDSMOpenInterestHistoryPath {
		t.Fatalf("request = %#v", executor.request)
	}
	wantQuery := map[string]string{"symbol": "BTCUSDT", "period": "5m", "limit": "500", "startTime": "1000", "endTime": "2000"}
	for key, want := range wantQuery {
		if got := executor.request.Query.Get(key); got != want {
			t.Errorf("query[%q] = %q, want %q", key, got, want)
		}
	}
	want := []OpenInterest{{
		Symbol: "BTCUSDT", SumOpenInterest: "20403.12345678", SumOpenInterestValue: "176196512.12345678",
		CMCCirculatingSupply: "165880.538", Timestamp: 1591261042378,
	}}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("result = %#v, want %#v", result, want)
	}
}

func TestSendOmitsOptionalParameters(t *testing.T) {
	executor := &recordingExecutor{body: []byte(`[]`)}
	client, _ := NewClient(executor)
	if _, err := client.Send(context.Background(), Params{Symbol: "BTCUSDT", Period: Period1h}); err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if executor.request.Query.Has("limit") || executor.request.Query.Has("startTime") || executor.request.Query.Has("endTime") {
		t.Fatalf("query = %v", executor.request.Query)
	}
}

func TestSendValidatesParameters(t *testing.T) {
	tests := []struct {
		name   string
		params Params
		want   string
	}{
		{name: "empty symbol", params: Params{Period: Period5m}, want: "symbol=empty"},
		{name: "invalid period", params: Params{Symbol: "BTCUSDT", Period: "10m"}, want: "period=invalid"},
		{name: "negative limit", params: Params{Symbol: "BTCUSDT", Period: Period5m, Limit: -1}, want: "limit=out_of_range"},
		{name: "large limit", params: Params{Symbol: "BTCUSDT", Period: Period5m, Limit: 501}, want: "limit=out_of_range"},
		{name: "negative start", params: Params{Symbol: "BTCUSDT", Period: Period5m, StartTime: -1}, want: "start_time=out_of_range"},
		{name: "negative end", params: Params{Symbol: "BTCUSDT", Period: Period5m, EndTime: -1}, want: "end_time=out_of_range"},
		{name: "reversed range", params: Params{Symbol: "BTCUSDT", Period: Period5m, StartTime: 2, EndTime: 1}, want: "time_range=invalid"},
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

func TestSendAcceptsSupportedPeriods(t *testing.T) {
	periods := []Period{Period5m, Period15m, Period30m, Period1h, Period2h, Period4h, Period6h, Period12h, Period1d}
	for _, period := range periods {
		client, _ := NewClient(&recordingExecutor{body: []byte(`[]`)})
		if _, err := client.Send(context.Background(), Params{Symbol: "BTCUSDT", Period: period}); err != nil {
			t.Errorf("Send() period %q error = %v", period, err)
		}
	}
}

func TestSendWrapsExecutorError(t *testing.T) {
	want := errors.New("transport failure")
	client, _ := NewClient(&recordingExecutor{err: want})
	_, err := client.Send(context.Background(), Params{Symbol: "BTCUSDT", Period: Period5m})
	if !errors.Is(err, want) {
		t.Fatalf("Send() error = %v, want wrapped error", err)
	}
}

func TestSendRejectsInvalidResponse(t *testing.T) {
	client, _ := NewClient(&recordingExecutor{body: []byte(`{`)})
	_, err := client.Send(context.Background(), Params{Symbol: "BTCUSDT", Period: Period5m})
	if err == nil || !strings.Contains(err.Error(), "failed to decode response body") {
		t.Fatalf("Send() error = %v", err)
	}
}

func TestSendRejectsNilClient(t *testing.T) {
	var client *Client
	_, err := client.Send(context.Background(), Params{Symbol: "BTCUSDT", Period: Period5m})
	if err == nil || !strings.Contains(err.Error(), "client=null") {
		t.Fatalf("Send() error = %v", err)
	}
}
