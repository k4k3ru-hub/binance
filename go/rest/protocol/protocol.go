// Package protocol contains wire-format types shared by Binance REST APIs.
package protocol

import (
	"encoding/json"
	"fmt"
)

// PriceLevel is one price and quantity pair in an order book.
type PriceLevel struct {
	Price    string
	Quantity string
}

func (l *PriceLevel) UnmarshalJSON(data []byte) error {
	var values []string
	if err := json.Unmarshal(data, &values); err != nil {
		return fmt.Errorf("failed to decode price level: %w", err)
	}
	if len(values) != 2 {
		return fmt.Errorf("failed to decode price level: expected 2 values, got %d", len(values))
	}
	l.Price = values[0]
	l.Quantity = values[1]
	return nil
}

func (l PriceLevel) MarshalJSON() ([]byte, error) {
	return json.Marshal([2]string{l.Price, l.Quantity})
}
