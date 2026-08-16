package exchange_info

import (
	"context"
	"fmt"
	"strings"

	"github.com/k4k3ru-hub/binance/go/rest"
	spotexchangeinfo "github.com/k4k3ru-hub/binance/go/rest/spot/exchange_info"
	"github.com/k4k3ru-hub/cli-go"
)

const (
	CommandName            = "exchangeInfo"
	OptionNameSymbol       = "symbol"
	OptionAliasSymbol      = "s"
	OptionNameSymbolStatus = "symbol-status"
)

// ParamsFromOptions converts CLI options into immutable API parameters.
func ParamsFromOptions(options map[string]*cli.Option) spotexchangeinfo.Params {
	params := spotexchangeinfo.Params{}
	if option := options[OptionNameSymbol]; option != nil {
		params.Symbol = strings.ToUpper(strings.TrimSpace(option.Value))
	}
	if option := options[OptionNameSymbolStatus]; option != nil {
		params.SymbolStatus = strings.ToUpper(strings.TrimSpace(option.Value))
	}
	return params
}

// Run executes the Spot exchange-information command.
func Run(options map[string]*cli.Option) {
	client, err := rest.NewClient(nil)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	result, err := client.Spot().ExchangeInfo().Send(context.Background(), ParamsFromOptions(options))
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	headers := []string{"Symbol", "Status", "BaseAsset", "QuoteAsset", "TickSize", "StepSize", "MinQty", "MinNotional", "OrderTypes"}
	data := make([][]interface{}, 0, len(result.Symbols))
	for _, symbol := range result.Symbols {
		var tickSize, stepSize, minQty, minNotional string
		for _, filter := range symbol.Filters {
			switch filter.FilterType {
			case "PRICE_FILTER":
				tickSize = filter.TickSize
			case "LOT_SIZE":
				stepSize, minQty = filter.StepSize, filter.MinQty
			case "MIN_NOTIONAL", "NOTIONAL":
				minNotional = filter.MinNotional
			}
		}
		data = append(data, []interface{}{
			symbol.Symbol, symbol.Status, symbol.BaseAsset, symbol.QuoteAsset,
			tickSize, stepSize, minQty, minNotional, strings.Join(symbol.OrderTypes, ","),
		})
	}
	cli.OutputTable(headers, data)
}

// SetCommand registers the command with its parent.
func SetCommand(parent *cli.Command) {
	command := cli.NewCommand(CommandName)
	command.Usage = "Get Spot exchange metadata."
	command.Options[OptionNameSymbol] = &cli.Option{Alias: OptionAliasSymbol, HasValue: true}
	command.Options[OptionNameSymbolStatus] = &cli.Option{HasValue: true}
	command.Action = Run
	parent.Commands = append(parent.Commands, command)
}
