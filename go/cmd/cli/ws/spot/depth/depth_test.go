package depth

import (
	"testing"

	spotdepth "github.com/k4k3ru-hub/binance/go/websocket/spot/depth"
	"github.com/k4k3ru-hub/cli-go"
)

func TestParamsFromOptions(t *testing.T) {
	params, err := ParamsFromOptions(map[string]*cli.Option{
		OptionNameSymbol: {Value: " btcusdt "}, OptionNameUpdateSpeed: {Value: "100MS"},
	})
	if err != nil {
		t.Fatalf("ParamsFromOptions() error = %v", err)
	}
	if params.Symbol != "BTCUSDT" || params.UpdateSpeed != spotdepth.UpdateSpeed100ms {
		t.Fatalf("params = %#v", params)
	}
}

func TestParamsFromOptionsRejectsInvalidInput(t *testing.T) {
	if _, err := ParamsFromOptions(nil); err == nil {
		t.Fatal("missing symbol error = nil")
	}
	if _, err := ParamsFromOptions(map[string]*cli.Option{OptionNameSymbol: {Value: "BTCUSDT"}, OptionNameUpdateSpeed: {Value: "500ms"}}); err == nil {
		t.Fatal("unsupported speed error = nil")
	}
}
