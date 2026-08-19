// Package usdsm composes Binance USDⓈ-M Futures REST operations.
package usdsm

import (
	"fmt"

	"github.com/k4k3ru-hub/binance/go/rest/transport"
	"github.com/k4k3ru-hub/binance/go/rest/usdsm/depth"
	"github.com/k4k3ru-hub/binance/go/rest/usdsm/exchange_info"
	fundinginfo "github.com/k4k3ru-hub/binance/go/rest/usdsm/funding_info"
	fundingratehistory "github.com/k4k3ru-hub/binance/go/rest/usdsm/funding_rate_history"
	openinterest "github.com/k4k3ru-hub/binance/go/rest/usdsm/open_interest"
	openinteresthistory "github.com/k4k3ru-hub/binance/go/rest/usdsm/open_interest_history"
	premiumindex "github.com/k4k3ru-hub/binance/go/rest/usdsm/premium_index"
)

type Client struct {
	exchangeInfo        *exchange_info.Client
	depth               *depth.Client
	openInterest        *openinterest.Client
	openInterestHistory *openinteresthistory.Client
	fundingRateHistory  *fundingratehistory.Client
	premiumIndex        *premiumindex.Client
	fundingInfo         *fundinginfo.Client
}

// NewClient creates a USDⓈ-M Futures API client and composes its operations.
//
// Parameters:
//   - executor: REST request executor.
//
// Returns:
//   - Composed USDⓈ-M Futures API client.
//
// Version:
//   - 2026-08-19: Added funding-rate operations.
//   - 2026-08-19: Added open-interest operations.
func NewClient(executor transport.Executor) (*Client, error) {
	if executor == nil {
		return nil, fmt.Errorf("failed to create USDⓈ-M API client: executor=null")
	}
	exchangeInfoClient, err := exchange_info.NewClient(executor)
	if err != nil {
		return nil, fmt.Errorf("failed to create USDⓈ-M API client: %w", err)
	}
	depthClient, err := depth.NewClient(executor)
	if err != nil {
		return nil, fmt.Errorf("failed to create USDⓈ-M API client: %w", err)
	}
	openInterestClient, err := openinterest.NewClient(executor)
	if err != nil {
		return nil, fmt.Errorf("failed to create USDⓈ-M API client: %w", err)
	}
	openInterestHistoryClient, err := openinteresthistory.NewClient(executor)
	if err != nil {
		return nil, fmt.Errorf("failed to create USDⓈ-M API client: %w", err)
	}
	fundingRateHistoryClient, err := fundingratehistory.NewClient(executor)
	if err != nil {
		return nil, fmt.Errorf("failed to create USDⓈ-M API client: %w", err)
	}
	premiumIndexClient, err := premiumindex.NewClient(executor)
	if err != nil {
		return nil, fmt.Errorf("failed to create USDⓈ-M API client: %w", err)
	}
	fundingInfoClient, err := fundinginfo.NewClient(executor)
	if err != nil {
		return nil, fmt.Errorf("failed to create USDⓈ-M API client: %w", err)
	}
	return &Client{
		exchangeInfo:        exchangeInfoClient,
		depth:               depthClient,
		openInterest:        openInterestClient,
		openInterestHistory: openInterestHistoryClient,
		fundingRateHistory:  fundingRateHistoryClient,
		premiumIndex:        premiumIndexClient,
		fundingInfo:         fundingInfoClient,
	}, nil
}

// ExchangeInfo returns the USDⓈ-M Futures exchange-information client.
//
// Returns:
//   - Exchange-information client.
//
// Version:
//   - 2026-08-19: Added documentation.
func (c *Client) ExchangeInfo() *exchange_info.Client {
	if c == nil {
		return nil
	}
	return c.exchangeInfo
}

// Depth returns the USDⓈ-M Futures order-book client.
//
// Returns:
//   - Order-book client.
//
// Version:
//   - 2026-08-19: Added documentation.
func (c *Client) Depth() *depth.Client {
	if c == nil {
		return nil
	}
	return c.depth
}

// OpenInterest returns the USDⓈ-M Futures current open-interest client.
//
// Returns:
//   - Current open-interest client.
//
// Version:
//   - 2026-08-19: Added.
func (c *Client) OpenInterest() *openinterest.Client {
	if c == nil {
		return nil
	}
	return c.openInterest
}

// OpenInterestHistory returns the USDⓈ-M Futures open-interest history client.
//
// Returns:
//   - Open-interest history client.
//
// Version:
//   - 2026-08-19: Added.
func (c *Client) OpenInterestHistory() *openinteresthistory.Client {
	if c == nil {
		return nil
	}
	return c.openInterestHistory
}

// FundingRateHistory returns the USDⓈ-M Futures funding-rate history client.
//
// Returns:
//   - Funding-rate history client.
//
// Version:
//   - 2026-08-19: Added.
func (c *Client) FundingRateHistory() *fundingratehistory.Client {
	if c == nil {
		return nil
	}
	return c.fundingRateHistory
}

// PremiumIndex returns the USDⓈ-M Futures mark-price and funding-rate client.
//
// Returns:
//   - Premium-index client.
//
// Version:
//   - 2026-08-19: Added.
func (c *Client) PremiumIndex() *premiumindex.Client {
	if c == nil {
		return nil
	}
	return c.premiumIndex
}

// FundingInfo returns the USDⓈ-M Futures funding-information client.
//
// Returns:
//   - Funding-information client.
//
// Version:
//   - 2026-08-19: Added.
func (c *Client) FundingInfo() *fundinginfo.Client {
	if c == nil {
		return nil
	}
	return c.fundingInfo
}
