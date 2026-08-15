package protocol

import (
	"encoding/json"
	"testing"
)

func TestDepthEventsDecodePriceLevels(t *testing.T) {
	var spot SpotDepthEvent
	if err := json.Unmarshal([]byte(`{"e":"depthUpdate","E":1,"s":"BTCUSDT","U":2,"u":3,"b":[["1","2"]],"a":[]}`), &spot); err != nil {
		t.Fatalf("SpotDepthEvent Unmarshal() error = %v", err)
	}
	if spot.Bids[0].Price != "1" || spot.Bids[0].Quantity != "2" {
		t.Fatalf("spot = %#v", spot)
	}

	var futures USDSMDepthEvent
	if err := json.Unmarshal([]byte(`{"e":"depthUpdate","E":1,"T":2,"s":"BTCUSDT","U":3,"u":4,"pu":2,"b":[],"a":[["5","6"]]}`), &futures); err != nil {
		t.Fatalf("USDSMDepthEvent Unmarshal() error = %v", err)
	}
	if futures.PreviousFinalUpdateID != 2 || futures.Asks[0].Quantity != "6" {
		t.Fatalf("futures = %#v", futures)
	}
}

func TestDecodeBookTickerEvents(t *testing.T) {
	spot, err := DecodeSpotBookTicker([]byte(`{"u":400900217,"s":"BNBUSDT","b":"25.35190000","B":"31.21000000","a":"25.36520000","A":"40.66000000"}`))
	if err != nil {
		t.Fatalf("DecodeSpotBookTicker() error = %v", err)
	}
	if spot.Symbol != "BNBUSDT" || spot.UpdateID != 400900217 || spot.BidQuantity != "31.21000000" || spot.AskPrice != "25.36520000" {
		t.Fatalf("spot = %#v", spot)
	}

	futures, err := DecodeUSDSMBookTicker([]byte(`{"stream":"btcusdt@bookTicker","data":{"e":"bookTicker","u":400900217,"E":1568014460893,"T":1568014460891,"s":"BTCUSDT","b":"25.35190000","B":"31.21000000","a":"25.36520000","A":"40.66000000"}}`))
	if err != nil {
		t.Fatalf("DecodeUSDSMBookTicker() error = %v", err)
	}
	if futures.EventType != "bookTicker" || futures.TransactionTime != 1568014460891 || futures.AskQuantity != "40.66000000" {
		t.Fatalf("futures = %#v", futures)
	}
}

func TestDecodeBookTickerRejectsMalformedMessages(t *testing.T) {
	if _, err := DecodeSpotBookTicker(nil); err == nil {
		t.Fatal("DecodeSpotBookTicker(nil) error = nil")
	}
	if _, err := DecodeUSDSMBookTicker([]byte(`{"stream":"btcusdt@bookTicker","data":null}`)); err == nil {
		t.Fatal("DecodeUSDSMBookTicker(empty data) error = nil")
	}
}
