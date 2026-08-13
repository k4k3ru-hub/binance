// Package binance provides the public facade for the Binance Go SDK.
package binance

import (
	"context"

	"github.com/k4k3ru-hub/binance/go/rest"
	"github.com/k4k3ru-hub/binance/go/rest/protocol"
	spotdepth "github.com/k4k3ru-hub/binance/go/rest/spot/depth"
	spotexchangeinfo "github.com/k4k3ru-hub/binance/go/rest/spot/exchange_info"
	usdsmdepth "github.com/k4k3ru-hub/binance/go/rest/usdsm/depth"
	usdsmexchangeinfo "github.com/k4k3ru-hub/binance/go/rest/usdsm/exchange_info"
	"github.com/k4k3ru-hub/binance/go/websocket"
	wsprotocol "github.com/k4k3ru-hub/binance/go/websocket/protocol"
	spotwsdepth "github.com/k4k3ru-hub/binance/go/websocket/spot/depth"
	usdsmwsdepth "github.com/k4k3ru-hub/binance/go/websocket/usdsm/depth"
)

type RESTClient = rest.Client
type RESTClientOption = rest.ClientOption
type RESTResponseError = rest.ResponseError

type SpotExchangeInfoParams = spotexchangeinfo.Params
type SpotExchangeInfo = spotexchangeinfo.ExchangeInfo
type SpotDepthParams = spotdepth.Params
type SpotDepth = spotdepth.Depth

type USDSMExchangeInfoParams = usdsmexchangeinfo.Params
type USDSMExchangeInfo = usdsmexchangeinfo.ExchangeInfo
type USDSMDepthParams = usdsmdepth.Params
type USDSMDepth = usdsmdepth.Depth

type PriceLevel = protocol.PriceLevel

type SpotWebSocketClient = websocket.SpotClient
type USDSMWebSocketClient = websocket.USDSMClient
type WebSocketClientOption = websocket.ClientOption
type WebSocketSessionHandler = websocket.SessionHandler
type WebSocketSessionContext = websocket.SessionContext

type SpotDepthSubscriptionParams = spotwsdepth.Params
type SpotDepthUpdateSpeed = spotwsdepth.UpdateSpeed
type USDSMDepthSubscriptionParams = usdsmwsdepth.Params
type USDSMDepthUpdateSpeed = usdsmwsdepth.UpdateSpeed
type SpotDepthEvent = wsprotocol.SpotDepthEvent
type USDSMDepthEvent = wsprotocol.USDSMDepthEvent

const (
	SpotDepthUpdateSpeedDefault  = spotwsdepth.UpdateSpeedDefault
	SpotDepthUpdateSpeed100ms    = spotwsdepth.UpdateSpeed100ms
	SpotDepthUpdateSpeed1000ms   = spotwsdepth.UpdateSpeed1000ms
	USDSMDepthUpdateSpeedDefault = usdsmwsdepth.UpdateSpeedDefault
	USDSMDepthUpdateSpeed100ms   = usdsmwsdepth.UpdateSpeed100ms
	USDSMDepthUpdateSpeed250ms   = usdsmwsdepth.UpdateSpeed250ms
	USDSMDepthUpdateSpeed500ms   = usdsmwsdepth.UpdateSpeed500ms
)

func NewRESTClient(option *RESTClientOption) (*RESTClient, error) { return rest.NewClient(option) }
func DefaultRESTClientOption() *RESTClientOption                  { return rest.DefaultClientOption() }

func NewSpotWebSocketClient(ctx context.Context, handler WebSocketSessionHandler, option *WebSocketClientOption) (*SpotWebSocketClient, error) {
	return websocket.NewSpotClient(ctx, handler, option)
}

func NewUSDSMWebSocketClient(ctx context.Context, handler WebSocketSessionHandler, option *WebSocketClientOption) (*USDSMWebSocketClient, error) {
	return websocket.NewUSDSMClient(ctx, handler, option)
}

func DefaultSpotWebSocketClientOption() *WebSocketClientOption {
	return websocket.DefaultSpotClientOption()
}
func DefaultUSDSMWebSocketClientOption() *WebSocketClientOption {
	return websocket.DefaultUSDSMClientOption()
}
