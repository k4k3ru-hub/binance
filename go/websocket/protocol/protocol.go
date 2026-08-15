// Package protocol defines Binance WebSocket wire-format types.
package protocol

import (
	"encoding/json"
	"fmt"
	"hash/fnv"

	"github.com/k4k3ru-hub/binance/go/rest/protocol"
)

const (
	MethodSubscribe   = "SUBSCRIBE"
	MethodUnsubscribe = "UNSUBSCRIBE"
)

// SubscriptionRequest is a live subscription control message.
type SubscriptionRequest struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	ID     int64    `json:"id"`
}

// SubscriptionResponse acknowledges a control message or reports an error.
type SubscriptionResponse struct {
	Result  json.RawMessage `json:"result,omitempty"`
	ID      int64           `json:"id,omitempty"`
	Code    int             `json:"code,omitempty"`
	Message string          `json:"msg,omitempty"`
}

// CombinedStreamEnvelope wraps an event received from a combined stream URL.
type CombinedStreamEnvelope struct {
	Stream string          `json:"stream"`
	Data   json.RawMessage `json:"data"`
}

// SpotDepthEvent is one Spot diff-depth update.
type SpotDepthEvent struct {
	EventType     string                `json:"e"`
	EventTime     int64                 `json:"E"`
	Symbol        string                `json:"s"`
	FirstUpdateID int64                 `json:"U"`
	FinalUpdateID int64                 `json:"u"`
	Bids          []protocol.PriceLevel `json:"b"`
	Asks          []protocol.PriceLevel `json:"a"`
}

// USDSMDepthEvent is one USDⓈ-M Futures diff-depth update.
type USDSMDepthEvent struct {
	EventType             string                `json:"e"`
	EventTime             int64                 `json:"E"`
	TransactionTime       int64                 `json:"T"`
	Symbol                string                `json:"s"`
	FirstUpdateID         int64                 `json:"U"`
	FinalUpdateID         int64                 `json:"u"`
	PreviousFinalUpdateID int64                 `json:"pu"`
	Bids                  []protocol.PriceLevel `json:"b"`
	Asks                  []protocol.PriceLevel `json:"a"`
}

// SpotBookTickerEvent is one Spot best bid/ask update.
type SpotBookTickerEvent struct {
	UpdateID    int64  `json:"u"`
	Symbol      string `json:"s"`
	BidPrice    string `json:"b"`
	BidQuantity string `json:"B"`
	AskPrice    string `json:"a"`
	AskQuantity string `json:"A"`
}

// USDSMBookTickerEvent is one USDⓈ-M Futures best bid/ask update.
type USDSMBookTickerEvent struct {
	EventType       string `json:"e"`
	EventTime       int64  `json:"E"`
	TransactionTime int64  `json:"T"`
	Symbol          string `json:"s"`
	UpdateID        int64  `json:"u"`
	BidPrice        string `json:"b"`
	BidQuantity     string `json:"B"`
	AskPrice        string `json:"a"`
	AskQuantity     string `json:"A"`
}

// DecodeSpotBookTicker decodes a raw or combined-stream Spot book-ticker message.
func DecodeSpotBookTicker(message []byte) (SpotBookTickerEvent, error) {
	var event SpotBookTickerEvent
	if err := decodeEvent(message, &event); err != nil {
		return event, fmt.Errorf("failed to decode Spot book ticker event: %w", err)
	}
	return event, nil
}

// DecodeUSDSMBookTicker decodes a raw or combined-stream USDⓈ-M book-ticker message.
func DecodeUSDSMBookTicker(message []byte) (USDSMBookTickerEvent, error) {
	var event USDSMBookTickerEvent
	if err := decodeEvent(message, &event); err != nil {
		return event, fmt.Errorf("failed to decode USDⓈ-M book ticker event: %w", err)
	}
	return event, nil
}

func decodeEvent(message []byte, event any) error {
	if len(message) == 0 {
		return fmt.Errorf("message=empty")
	}
	var envelope CombinedStreamEnvelope
	if err := json.Unmarshal(message, &envelope); err != nil {
		return err
	}
	if envelope.Stream != "" {
		if len(envelope.Data) == 0 || string(envelope.Data) == "null" {
			return fmt.Errorf("combined stream data=empty")
		}
		return json.Unmarshal(envelope.Data, event)
	}
	return json.Unmarshal(message, event)
}

// RequestID returns a stable positive request identifier for one operation.
func RequestID(method, stream string) int64 {
	hash := fnv.New64a()
	_, _ = hash.Write([]byte(method))
	_, _ = hash.Write([]byte{0})
	_, _ = hash.Write([]byte(stream))
	return int64(hash.Sum64() & 0x7fffffffffffffff)
}
