# Binance SDK

A Go SDK for selected Binance Spot and USDⓈ-M Futures REST and WebSocket APIs.

The Go module is located in [`go/`](go/README.md).

It includes a CLI for all implemented endpoints under `go/cmd/cli`.

## Implemented APIs

- Spot `GET /api/v3/exchangeInfo`
- Spot `GET /api/v3/depth`
- USDⓈ-M Futures `GET /fapi/v1/exchangeInfo`
- USDⓈ-M Futures `GET /fapi/v1/depth`
- Spot WebSocket diff-depth stream
- USDⓈ-M Futures WebSocket diff-depth stream
- Spot WebSocket trade stream
- USDⓈ-M Futures WebSocket aggregate-trade stream

## References

- [Binance Spot REST API](https://developers.binance.com/docs/binance-spot-api-docs/rest-api)
- [Binance USDⓈ-M Futures REST API](https://developers.binance.com/docs/derivatives/usds-margined-futures/market-data/rest-api)
