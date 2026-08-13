package exchange_info

import (
	"testing"

	"github.com/k4k3ru-hub/cli-go"
)

func TestParamsFromOptionsNormalizesValues(t *testing.T) {
	params := ParamsFromOptions(map[string]*cli.Option{
		OptionNameSymbol:       {Value: " btcusdt "},
		OptionNameSymbolStatus: {Value: "trading"},
	})
	if params.Symbol != "BTCUSDT" || params.SymbolStatus != "TRADING" {
		t.Fatalf("ParamsFromOptions() = %#v", params)
	}
}
