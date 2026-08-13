// Package exchange_info implements the Binance USDⓈ-M Futures exchangeInfo endpoint.
package exchange_info

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/k4k3ru-hub/binance/go/rest/endpoint"
	"github.com/k4k3ru-hub/binance/go/rest/transport"
)

type Client struct{ executor transport.Executor }
type Params struct{}

type ExchangeInfo struct {
	Timezone        string            `json:"timezone"`
	ServerTime      int64             `json:"serverTime"`
	FuturesType     string            `json:"futuresType"`
	RateLimits      []RateLimit       `json:"rateLimits"`
	ExchangeFilters []json.RawMessage `json:"exchangeFilters"`
	Assets          []Asset           `json:"assets"`
	Symbols         []Symbol          `json:"symbols"`
}

type RateLimit struct {
	RateLimitType string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	IntervalNum   int    `json:"intervalNum"`
	Limit         int    `json:"limit"`
}

type Asset struct {
	Asset             string `json:"asset"`
	MarginAvailable   bool   `json:"marginAvailable"`
	AutoAssetExchange string `json:"autoAssetExchange"`
}

type Symbol struct {
	Symbol                string     `json:"symbol"`
	Pair                  string     `json:"pair"`
	ContractType          string     `json:"contractType"`
	DeliveryDate          int64      `json:"deliveryDate"`
	OnboardDate           int64      `json:"onboardDate"`
	Status                string     `json:"status"`
	MaintMarginPercent    string     `json:"maintMarginPercent"`
	RequiredMarginPercent string     `json:"requiredMarginPercent"`
	BaseAsset             string     `json:"baseAsset"`
	QuoteAsset            string     `json:"quoteAsset"`
	MarginAsset           string     `json:"marginAsset"`
	PricePrecision        int        `json:"pricePrecision"`
	QuantityPrecision     int        `json:"quantityPrecision"`
	BaseAssetPrecision    int        `json:"baseAssetPrecision"`
	QuotePrecision        int        `json:"quotePrecision"`
	UnderlyingType        string     `json:"underlyingType"`
	UnderlyingSubType     []string   `json:"underlyingSubType"`
	SettlePlan            int        `json:"settlePlan"`
	TriggerProtect        string     `json:"triggerProtect"`
	LiquidationFee        string     `json:"liquidationFee"`
	MarketTakeBound       string     `json:"marketTakeBound"`
	MaxMoveOrderLimit     int        `json:"maxMoveOrderLimit"`
	Filters               []Filter   `json:"filters"`
	OrderTypes            []string   `json:"orderTypes"`
	TimeInForce           []string   `json:"timeInForce"`
	PermissionSets        [][]string `json:"permissionSets,omitempty"`
}

// UnmarshalJSON decodes a USDⓈ-M symbol and normalizes permission sets.
//
// Notes:
//   - Binance may return permissionSets as either []string or [][]string.
//   - A one-dimensional value is normalized to one permission set.
//
// Parameters:
//   - data: Encoded USDⓈ-M symbol.
//
// Returns:
//   - Decode error.
//
// Version:
//   - 2026-08-14: Accepted one-dimensional permission sets.
func (s *Symbol) UnmarshalJSON(data []byte) error {
	type symbolAlias Symbol
	var raw struct {
		symbolAlias
		PermissionSets json.RawMessage `json:"permissionSets"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("failed to decode USDⓈ-M symbol: %w", err)
	}

	*s = Symbol(raw.symbolAlias)
	if len(raw.PermissionSets) == 0 || string(raw.PermissionSets) == "null" {
		return nil
	}

	var nested [][]string
	if err := json.Unmarshal(raw.PermissionSets, &nested); err == nil {
		s.PermissionSets = nested
		return nil
	}

	var flat []string
	if err := json.Unmarshal(raw.PermissionSets, &flat); err != nil {
		return fmt.Errorf("failed to decode USDⓈ-M symbol: failed to decode permission sets: %w", err)
	}
	if len(flat) != 0 {
		s.PermissionSets = [][]string{flat}
	}

	return nil
}

// Filter represents the union of documented USDⓈ-M symbol-filter fields.
type Filter struct {
	FilterType        string `json:"filterType"`
	MinPrice          string `json:"minPrice,omitempty"`
	MaxPrice          string `json:"maxPrice,omitempty"`
	TickSize          string `json:"tickSize,omitempty"`
	MinQty            string `json:"minQty,omitempty"`
	MaxQty            string `json:"maxQty,omitempty"`
	StepSize          string `json:"stepSize,omitempty"`
	Limit             int    `json:"limit,omitempty"`
	Notional          string `json:"notional,omitempty"`
	MultiplierUp      string `json:"multiplierUp,omitempty"`
	MultiplierDown    string `json:"multiplierDown,omitempty"`
	MultiplierDecimal string `json:"multiplierDecimal,omitempty"`
}

func NewClient(executor transport.Executor) (*Client, error) {
	if executor == nil {
		return nil, fmt.Errorf("failed to create USDⓈ-M exchange-info client: executor=null")
	}
	return &Client{executor: executor}, nil
}

func (c *Client) Send(ctx context.Context, _ Params) (*ExchangeInfo, error) {
	if c == nil {
		return nil, fmt.Errorf("failed to request USDⓈ-M exchange info: client=null")
	}
	body, err := c.executor.Do(ctx, transport.Request{Method: http.MethodGet, Path: endpoint.USDSMExchangeInfoPath})
	if err != nil {
		return nil, fmt.Errorf("failed to request USDⓈ-M exchange info: %w", err)
	}
	var result ExchangeInfo
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to request USDⓈ-M exchange info: failed to decode response body: %w", err)
	}
	return &result, nil
}
