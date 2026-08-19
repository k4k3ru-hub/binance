// Package endpoint defines Binance REST endpoints.
package endpoint

const (
	DefaultSpotBaseURL  = "https://api.binance.com"
	DefaultUSDSMBaseURL = "https://fapi.binance.com"

	SpotExchangeInfoPath         = "/api/v3/exchangeInfo"
	SpotDepthPath                = "/api/v3/depth"
	USDSMExchangeInfoPath        = "/fapi/v1/exchangeInfo"
	USDSMDepthPath               = "/fapi/v1/depth"
	USDSMOpenInterestPath        = "/fapi/v1/openInterest"
	USDSMOpenInterestHistoryPath = "/futures/data/openInterestHist"
	USDSMFundingRateHistoryPath  = "/fapi/v1/fundingRate"
	USDSMPremiumIndexPath        = "/fapi/v1/premiumIndex"
	USDSMFundingInfoPath         = "/fapi/v1/fundingInfo"
)
