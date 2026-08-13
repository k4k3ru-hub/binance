// Package protocol defines Binance WebSocket wire-format types.
package protocol

import (
	"encoding/json"
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

// RequestID returns a stable positive request identifier for one operation.
func RequestID(method, stream string) int64 {
	hash := fnv.New64a()
	_, _ = hash.Write([]byte(method))
	_, _ = hash.Write([]byte{0})
	_, _ = hash.Write([]byte(stream))
	return int64(hash.Sum64() & 0x7fffffffffffffff)
}
