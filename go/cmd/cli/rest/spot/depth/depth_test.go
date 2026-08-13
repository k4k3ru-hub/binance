package depth

import (
	"testing"

	"github.com/k4k3ru-hub/cli-go"
)

func TestParamsFromOptions(t *testing.T) {
	params, err := ParamsFromOptions(map[string]*cli.Option{
		OptionNameSymbol: {Value: " btcusdt "}, OptionNameLimit: {Value: "100"},
	})
	if err != nil {
		t.Fatalf("ParamsFromOptions() error = %v", err)
	}
	if params.Symbol != "BTCUSDT" || params.Limit != 100 {
		t.Fatalf("ParamsFromOptions() = %#v", params)
	}
}

func TestParamsFromOptionsRejectsInvalidInput(t *testing.T) {
	if _, err := ParamsFromOptions(nil); err == nil {
		t.Fatal("missing symbol error = nil")
	}
	if _, err := ParamsFromOptions(map[string]*cli.Option{OptionNameSymbol: {Value: "BTCUSDT"}, OptionNameLimit: {Value: "x"}}); err == nil {
		t.Fatal("invalid limit error = nil")
	}
}
