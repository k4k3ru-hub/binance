package depth

import (
	"testing"

	"github.com/k4k3ru-hub/cli-go"
)

func TestParamsFromOptions(t *testing.T) {
	params, err := ParamsFromOptions(map[string]*cli.Option{
		OptionNameSymbol: {Value: " ethusdt "}, OptionNameLimit: {Value: "500"},
	})
	if err != nil {
		t.Fatalf("ParamsFromOptions() error = %v", err)
	}
	if params.Symbol != "ETHUSDT" || params.Limit != 500 {
		t.Fatalf("ParamsFromOptions() = %#v", params)
	}
}

func TestParamsFromOptionsRejectsInvalidInput(t *testing.T) {
	if _, err := ParamsFromOptions(nil); err == nil {
		t.Fatal("missing symbol error = nil")
	}
	if _, err := ParamsFromOptions(map[string]*cli.Option{OptionNameSymbol: {Value: "BTCUSDT"}, OptionNameLimit: {Value: "0"}}); err == nil {
		t.Fatal("invalid limit error = nil")
	}
}
