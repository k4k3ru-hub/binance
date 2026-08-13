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
