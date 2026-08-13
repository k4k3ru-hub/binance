package depth

import (
	"testing"

	usdsmdepth "github.com/k4k3ru-hub/binance/go/websocket/usdsm/depth"
	"github.com/k4k3ru-hub/cli-go"
)

func TestParamsFromOptions(t *testing.T) {
	params, err := ParamsFromOptions(map[string]*cli.Option{
		OptionNameSymbol: {Value: " ethusdt "}, OptionNameUpdateSpeed: {Value: "500MS"},
	})
	if err != nil {
		t.Fatalf("ParamsFromOptions() error = %v", err)
	}
	if params.Symbol != "ETHUSDT" || params.UpdateSpeed != usdsmdepth.UpdateSpeed500ms {
		t.Fatalf("params = %#v", params)
	}
}

func TestParamsFromOptionsRejectsInvalidSpeed(t *testing.T) {
	if _, err := ParamsFromOptions(map[string]*cli.Option{OptionNameSymbol: {Value: "BTCUSDT"}, OptionNameUpdateSpeed: {Value: "1000ms"}}); err == nil {
		t.Fatal("unsupported speed error = nil")
	}
}
