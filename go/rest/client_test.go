package rest

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	spotdepth "github.com/k4k3ru-hub/binance/go/rest/spot/depth"
	usdsmdepth "github.com/k4k3ru-hub/binance/go/rest/usdsm/depth"
)

type routingHTTPClient struct{ urls []string }

func (c *routingHTTPClient) Do(request *http.Request) (*http.Response, error) {
	c.urls = append(c.urls, request.URL.String())
	body := `{"lastUpdateId":1,"bids":[],"asks":[]}`
	return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(body))}, nil
}

func TestClientRoutesAPIGroupsToSeparateBaseURLs(t *testing.T) {
	httpClient := &routingHTTPClient{}
	client, err := NewClient(&ClientOption{
		SpotBaseURL: "https://spot.example", USDSMBaseURL: "https://futures.example", HTTPClient: httpClient,
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if _, err := client.Spot().Depth().Send(context.Background(), spotdepth.Params{Symbol: "BTCUSDT"}); err != nil {
		t.Fatalf("Spot Depth Send() error = %v", err)
	}
	if _, err := client.USDSM().Depth().Send(context.Background(), usdsmdepth.Params{Symbol: "BTCUSDT"}); err != nil {
		t.Fatalf("USDSM Depth Send() error = %v", err)
	}
	if len(httpClient.urls) != 2 || httpClient.urls[0] != "https://spot.example/api/v3/depth?symbol=BTCUSDT" || httpClient.urls[1] != "https://futures.example/fapi/v1/depth?symbol=BTCUSDT" {
		t.Fatalf("requested URLs = %#v", httpClient.urls)
	}
}

func TestClientComposesUSDSMOpenInterestOperations(t *testing.T) {
	client, err := NewClient(&ClientOption{HTTPClient: &routingHTTPClient{}})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if client.USDSM().OpenInterest() == nil {
		t.Fatal("USDSM OpenInterest() = nil")
	}
	if client.USDSM().OpenInterestHistory() == nil {
		t.Fatal("USDSM OpenInterestHistory() = nil")
	}
}

func TestClientComposesUSDSMFundingOperations(t *testing.T) {
	client, err := NewClient(&ClientOption{HTTPClient: &routingHTTPClient{}})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if client.USDSM().FundingRateHistory() == nil {
		t.Fatal("USDSM FundingRateHistory() = nil")
	}
	if client.USDSM().PremiumIndex() == nil {
		t.Fatal("USDSM PremiumIndex() = nil")
	}
	if client.USDSM().FundingInfo() == nil {
		t.Fatal("USDSM FundingInfo() = nil")
	}
}

type errorHTTPClient struct{}

func (errorHTTPClient) Do(*http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: http.StatusBadRequest, Body: io.NopCloser(strings.NewReader(`{"code":-1121,"msg":"Invalid symbol."}`))}, nil
}

func TestResponseErrorDecodesBinanceError(t *testing.T) {
	client, err := NewClient(&ClientOption{HTTPClient: errorHTTPClient{}})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	_, err = client.Spot().Depth().Send(context.Background(), spotdepth.Params{Symbol: "INVALID"})
	var responseError *ResponseError
	if !errors.As(err, &responseError) {
		t.Fatalf("Send() error = %v, want ResponseError", err)
	}
	if responseError.StatusCode != http.StatusBadRequest || responseError.Code != -1121 || !strings.Contains(err.Error(), "Invalid symbol") {
		t.Fatalf("ResponseError = %#v", responseError)
	}
}
