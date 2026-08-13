// Package exchange_info implements the Binance Spot exchangeInfo endpoint.
package exchange_info

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/k4k3ru-hub/binance/go/rest/endpoint"
	"github.com/k4k3ru-hub/binance/go/rest/transport"
)

type Client struct{ executor transport.Executor }

// Params configures a Spot exchange-information request. All fields are optional.
type Params struct {
	Symbol             string
	Symbols            []string
	Permissions        []string
	ShowPermissionSets *bool
	SymbolStatus       string
}

type ExchangeInfo struct {
	Timezone        string            `json:"timezone"`
	ServerTime      int64             `json:"serverTime"`
	RateLimits      []RateLimit       `json:"rateLimits"`
	ExchangeFilters []json.RawMessage `json:"exchangeFilters"`
	Symbols         []Symbol          `json:"symbols"`
	SORs            []SOR             `json:"sors,omitempty"`
}

type RateLimit struct {
	RateLimitType string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	IntervalNum   int    `json:"intervalNum"`
	Limit         int    `json:"limit"`
}

type Symbol struct {
	Symbol                          string     `json:"symbol"`
	Status                          string     `json:"status"`
	BaseAsset                       string     `json:"baseAsset"`
	BaseAssetPrecision              int        `json:"baseAssetPrecision"`
	QuoteAsset                      string     `json:"quoteAsset"`
	QuotePrecision                  int        `json:"quotePrecision"`
	QuoteAssetPrecision             int        `json:"quoteAssetPrecision"`
	BaseCommissionPrecision         int        `json:"baseCommissionPrecision"`
	QuoteCommissionPrecision        int        `json:"quoteCommissionPrecision"`
	OrderTypes                      []string   `json:"orderTypes"`
	IcebergAllowed                  bool       `json:"icebergAllowed"`
	OCOAllowed                      bool       `json:"ocoAllowed"`
	OTOAllowed                      bool       `json:"otoAllowed"`
	OPOAllowed                      bool       `json:"opoAllowed"`
	QuoteOrderQtyMarketAllowed      bool       `json:"quoteOrderQtyMarketAllowed"`
	AllowTrailingStop               bool       `json:"allowTrailingStop"`
	CancelReplaceAllowed            bool       `json:"cancelReplaceAllowed"`
	AmendAllowed                    bool       `json:"amendAllowed"`
	PegInstructionsAllowed          bool       `json:"pegInstructionsAllowed"`
	IsSpotTradingAllowed            bool       `json:"isSpotTradingAllowed"`
	IsMarginTradingAllowed          bool       `json:"isMarginTradingAllowed"`
	Filters                         []Filter   `json:"filters"`
	Permissions                     []string   `json:"permissions"`
	PermissionSets                  [][]string `json:"permissionSets"`
	DefaultSelfTradePreventionMode  string     `json:"defaultSelfTradePreventionMode"`
	AllowedSelfTradePreventionModes []string   `json:"allowedSelfTradePreventionModes"`
}

// Filter represents the union of documented Spot symbol-filter fields.
type Filter struct {
	FilterType            string `json:"filterType"`
	MinPrice              string `json:"minPrice,omitempty"`
	MaxPrice              string `json:"maxPrice,omitempty"`
	TickSize              string `json:"tickSize,omitempty"`
	MultiplierUp          string `json:"multiplierUp,omitempty"`
	MultiplierDown        string `json:"multiplierDown,omitempty"`
	BidMultiplierUp       string `json:"bidMultiplierUp,omitempty"`
	BidMultiplierDown     string `json:"bidMultiplierDown,omitempty"`
	AskMultiplierUp       string `json:"askMultiplierUp,omitempty"`
	AskMultiplierDown     string `json:"askMultiplierDown,omitempty"`
	AvgPriceMins          int    `json:"avgPriceMins,omitempty"`
	MinQty                string `json:"minQty,omitempty"`
	MaxQty                string `json:"maxQty,omitempty"`
	StepSize              string `json:"stepSize,omitempty"`
	MinNotional           string `json:"minNotional,omitempty"`
	MaxNotional           string `json:"maxNotional,omitempty"`
	ApplyToMarket         bool   `json:"applyToMarket,omitempty"`
	ApplyMinToMarket      bool   `json:"applyMinToMarket,omitempty"`
	ApplyMaxToMarket      bool   `json:"applyMaxToMarket,omitempty"`
	Limit                 int    `json:"limit,omitempty"`
	MaxNumOrders          int    `json:"maxNumOrders,omitempty"`
	MaxNumAlgoOrders      int    `json:"maxNumAlgoOrders,omitempty"`
	MaxNumIcebergOrders   int    `json:"maxNumIcebergOrders,omitempty"`
	MaxNumOrderAmends     int    `json:"maxNumOrderAmends,omitempty"`
	MaxNumOrderLists      int    `json:"maxNumOrderLists,omitempty"`
	MaxPosition           string `json:"maxPosition,omitempty"`
	MinTrailingAboveDelta int    `json:"minTrailingAboveDelta,omitempty"`
	MaxTrailingAboveDelta int    `json:"maxTrailingAboveDelta,omitempty"`
	MinTrailingBelowDelta int    `json:"minTrailingBelowDelta,omitempty"`
	MaxTrailingBelowDelta int    `json:"maxTrailingBelowDelta,omitempty"`
}

type SOR struct {
	BaseAsset string   `json:"baseAsset"`
	Symbols   []string `json:"symbols"`
}

func NewClient(executor transport.Executor) (*Client, error) {
	if executor == nil {
		return nil, fmt.Errorf("failed to create Spot exchange-info client: executor=null")
	}
	return &Client{executor: executor}, nil
}

func (c *Client) Send(ctx context.Context, params Params) (*ExchangeInfo, error) {
	if c == nil {
		return nil, fmt.Errorf("failed to request Spot exchange info: client=null")
	}
	query := url.Values{}
	if params.Symbol != "" {
		query.Set("symbol", params.Symbol)
	}
	if len(params.Symbols) != 0 {
		value, err := json.Marshal(params.Symbols)
		if err != nil {
			return nil, fmt.Errorf("failed to request Spot exchange info: failed to encode symbols: %w", err)
		}
		query.Set("symbols", string(value))
	}
	if len(params.Permissions) != 0 {
		value, err := json.Marshal(params.Permissions)
		if err != nil {
			return nil, fmt.Errorf("failed to request Spot exchange info: failed to encode permissions: %w", err)
		}
		query.Set("permissions", string(value))
	}
	if params.ShowPermissionSets != nil {
		query.Set("showPermissionSets", strconv.FormatBool(*params.ShowPermissionSets))
	}
	if params.SymbolStatus != "" {
		query.Set("symbolStatus", params.SymbolStatus)
	}
	body, err := c.executor.Do(ctx, transport.Request{Method: http.MethodGet, Path: endpoint.SpotExchangeInfoPath, Query: query})
	if err != nil {
		return nil, fmt.Errorf("failed to request Spot exchange info: %w", err)
	}
	var result ExchangeInfo
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to request Spot exchange info: failed to decode response body: %w", err)
	}
	return &result, nil
}
