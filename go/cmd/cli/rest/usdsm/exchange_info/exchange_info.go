package exchange_info

import (
	"context"
	"fmt"
	"strings"

	binance "github.com/k4k3ru-hub/binance/go"
	usdsmexchangeinfo "github.com/k4k3ru-hub/binance/go/rest/usdsm/exchange_info"
	"github.com/k4k3ru-hub/cli-go"
)

const CommandName = "exchangeInfo"

// Run executes the USDⓈ-M Futures exchange-information command.
func Run(_ map[string]*cli.Option) {
	client, err := binance.NewRESTClient(nil)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	result, err := client.USDSM().ExchangeInfo().Send(context.Background(), usdsmexchangeinfo.Params{})
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	headers := []string{"Symbol", "Status", "ContractType", "BaseAsset", "QuoteAsset", "MarginAsset", "TickSize", "StepSize", "MinQty", "Notional", "OrderTypes"}
	data := make([][]interface{}, 0, len(result.Symbols))
	for _, symbol := range result.Symbols {
		var tickSize, stepSize, minQty, notional string
		for _, filter := range symbol.Filters {
			switch filter.FilterType {
			case "PRICE_FILTER":
				tickSize = filter.TickSize
			case "LOT_SIZE":
				stepSize, minQty = filter.StepSize, filter.MinQty
			case "MIN_NOTIONAL":
				notional = filter.Notional
			}
		}
		data = append(data, []interface{}{
			symbol.Symbol, symbol.Status, symbol.ContractType, symbol.BaseAsset, symbol.QuoteAsset,
			symbol.MarginAsset, tickSize, stepSize, minQty, notional, strings.Join(symbol.OrderTypes, ","),
		})
	}
	cli.OutputTable(headers, data)
}

// SetCommand registers the command with its parent.
func SetCommand(parent *cli.Command) {
	command := cli.NewCommand(CommandName)
	command.Usage = "Get USDⓈ-M Futures exchange metadata."
	command.Action = Run
	parent.Commands = append(parent.Commands, command)
}
