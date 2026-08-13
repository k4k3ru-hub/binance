package binance

import (
	"errors"
	"net/http"
	"testing"
)

type failingHTTPClient struct{}

func (failingHTTPClient) Do(*http.Request) (*http.Response, error) {
	return nil, errors.New("unexpected request")
}

func TestNewRESTClientComposesAllAPIGroups(t *testing.T) {
	client, err := NewRESTClient(&RESTClientOption{HTTPClient: failingHTTPClient{}})
	if err != nil {
		t.Fatalf("NewRESTClient() error = %v", err)
	}
	if client.Spot() == nil || client.Spot().ExchangeInfo() == nil || client.Spot().Depth() == nil {
		t.Fatal("NewRESTClient() did not compose the Spot facade")
	}
	if client.USDSM() == nil || client.USDSM().ExchangeInfo() == nil || client.USDSM().Depth() == nil {
		t.Fatal("NewRESTClient() did not compose the USDⓈ-M facade")
	}
}

func TestDefaultRESTClientOptionHasBothBaseURLs(t *testing.T) {
	option := DefaultRESTClientOption()
	if option == nil || option.SpotBaseURL == "" || option.USDSMBaseURL == "" {
		t.Fatalf("DefaultRESTClientOption() = %#v", option)
	}
}

type testWebSocketHandler struct{}

func (*testWebSocketHandler) HandleMessage(WebSocketSessionContext, []byte) {}
func (*testWebSocketHandler) HandleClose(WebSocketSessionContext)           {}

func TestWebSocketFacadeConstructorsComposeDepthClients(t *testing.T) {
	spot, err := NewSpotWebSocketClient(nil, &testWebSocketHandler{}, nil)
	if err != nil || spot == nil || spot.Depth() == nil {
		t.Fatalf("NewSpotWebSocketClient() = %#v, %v", spot, err)
	}
	futures, err := NewUSDSMWebSocketClient(nil, &testWebSocketHandler{}, nil)
	if err != nil || futures == nil || futures.Depth() == nil {
		t.Fatalf("NewUSDSMWebSocketClient() = %#v, %v", futures, err)
	}
}
