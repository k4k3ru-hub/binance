# Binance SDK for Go

This module provides a root facade over Binance Spot and USDⓈ-M Futures market-data APIs. The two API groups use separate Binance base URLs while sharing an injectable HTTP client.

## Requirements

- Go 1.25.4 or later

## CLI

Build the command:

```sh
go build -o build/binance-cli ./cmd/cli
```

Available commands:

```sh
./build/binance-cli rest spot exchangeInfo
./build/binance-cli rest spot exchangeInfo --symbol BTCUSDT
./build/binance-cli rest spot depth --symbol BTCUSDT --limit 100
./build/binance-cli rest usdsm exchangeInfo
./build/binance-cli rest usdsm depth --symbol BTCUSDT --limit 100

./build/binance-cli ws spot depth --symbol BTCUSDT --update-speed 100ms
./build/binance-cli ws usdsm depth --symbol BTCUSDT --update-speed 500ms
```

`--symbol`, `--limit`, and `--update-speed` have the short aliases `-s`, `-l`, and `-u`.

## Usage

Create the composed REST client:

```go
client, err := binance.NewRESTClient(nil)
if err != nil {
	return err
}
```

Get Spot exchange metadata:

```go
info, err := client.Spot().ExchangeInfo().Send(ctx, binance.SpotExchangeInfoParams{
	Symbol: "BTCUSDT",
})
```

Get a Spot order-book snapshot:

```go
book, err := client.Spot().Depth().Send(ctx, binance.SpotDepthParams{
	Symbol: "BTCUSDT",
	Limit:  100,
})
```

Get USDⓈ-M Futures metadata and an order-book snapshot:

```go
info, err := client.USDSM().ExchangeInfo().Send(ctx, binance.USDSMExchangeInfoParams{})

book, err := client.USDSM().Depth().Send(ctx, binance.USDSMDepthParams{
	Symbol: "BTCUSDT",
	Limit:  100,
})
```

The complete examples require:

```go
import binance "github.com/k4k3ru-hub/binance/go"
```

## Custom configuration

```go
client, err := binance.NewRESTClient(&binance.RESTClientOption{
	SpotBaseURL:  "http://127.0.0.1:8080",
	USDSMBaseURL: "http://127.0.0.1:8081",
	HTTPClient:   customHTTPClient,
})
```

## WebSocket SDK

Import the WebSocket client, protocol, and operation packages directly:

```go
import (
	"github.com/k4k3ru-hub/binance/go/websocket"
	"github.com/k4k3ru-hub/binance/go/websocket/protocol"
	spotbookticker "github.com/k4k3ru-hub/binance/go/websocket/spot/book_ticker"
)
```

Implement a session handler and decode book-ticker messages:

```go
type handler struct{}

func (*handler) HandleMessage(_ websocket.SessionContext, message []byte) {
	event, err := protocol.DecodeSpotBookTicker(message)
	if err != nil {
		return
	}
	// Process event.BidPrice, event.BidQuantity, event.AskPrice, and event.AskQuantity.
}

func (*handler) HandleClose(websocket.SessionContext) {}
```

Create a market-specific client and subscribe:

```go
client, err := websocket.NewSpotClient(ctx, &handler{}, nil)
if err != nil {
	return err
}
defer client.Close()

err = client.BookTicker().Subscribe(ctx, spotbookticker.Params{
	Symbol: "BTCUSDT",
})
```

Call `BookTicker().Unsubscribe` with the same parameters to stop the stream. For USDⓈ-M Futures, import `websocket/usdsm/book_ticker`, call `websocket.NewUSDSMClient`, and decode events with `protocol.DecodeUSDSMBookTicker`. Subscriptions are retained for restoration after reconnects.

Depth events are incremental updates. Build a consistent local order book by buffering WebSocket events and combining them with the corresponding REST depth snapshot.

## Package layout

```text
facade.go               Public facade and common type aliases
cmd/cli/                CLI entry point and REST/WebSocket subcommands
rest/
├── endpoint/           REST endpoints and default base URLs
├── protocol/           Shared wire-format types
├── transport/          Stateless request execution contract
├── spot/               Spot composition and operations
└── usdsm/              USDⓈ-M Futures composition and operations
websocket/
├── protocol/           Subscription messages, event types, and decoders
├── subscriptions/      Subscription execution contract
├── spot/               Spot depth and book-ticker subscriptions
└── usdsm/              USDⓈ-M Futures depth and book-ticker subscriptions
```
