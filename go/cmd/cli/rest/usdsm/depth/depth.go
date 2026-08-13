package depth

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	binance "github.com/k4k3ru-hub/binance/go"
	usdsmdepth "github.com/k4k3ru-hub/binance/go/rest/usdsm/depth"
	"github.com/k4k3ru-hub/cli-go"
)

const (
	CommandName       = "depth"
	OptionNameSymbol  = "symbol"
	OptionAliasSymbol = "s"
	OptionNameLimit   = "limit"
	OptionAliasLimit  = "l"
)

var command = cli.NewCommand(CommandName)

// ParamsFromOptions validates CLI options and returns immutable API parameters.
func ParamsFromOptions(options map[string]*cli.Option) (usdsmdepth.Params, error) {
	params := usdsmdepth.Params{}
	if option := options[OptionNameSymbol]; option != nil {
		params.Symbol = strings.ToUpper(strings.TrimSpace(option.Value))
	}
	if params.Symbol == "" {
		return params, fmt.Errorf("symbol is required")
	}
	if option := options[OptionNameLimit]; option != nil && strings.TrimSpace(option.Value) != "" {
		limit, err := strconv.Atoi(option.Value)
		if err != nil || limit <= 0 {
			return params, fmt.Errorf("limit must be a positive integer")
		}
		params.Limit = limit
	}
	return params, nil
}

// Run executes the USDⓈ-M Futures depth command.
func Run(options map[string]*cli.Option) {
	params, err := ParamsFromOptions(options)
	if err != nil {
		fmt.Printf("%s\n", err)
		command.ShowUsage()
		return
	}
	client, err := binance.NewRESTClient(nil)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	result, err := client.USDSM().Depth().Send(context.Background(), params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	headers := []string{"Side", "Price", "Quantity"}
	data := make([][]interface{}, 0, len(result.Bids)+len(result.Asks))
	for _, level := range result.Bids {
		data = append(data, []interface{}{"BID", level.Price, level.Quantity})
	}
	for _, level := range result.Asks {
		data = append(data, []interface{}{"ASK", level.Price, level.Quantity})
	}
	cli.OutputTable(headers, data)
}

// SetCommand registers the command with its parent.
func SetCommand(parent *cli.Command) {
	command.Usage = "Get a USDⓈ-M Futures order-book snapshot."
	command.Options[OptionNameSymbol] = &cli.Option{Alias: OptionAliasSymbol, HasValue: true}
	command.Options[OptionNameLimit] = &cli.Option{Alias: OptionAliasLimit, HasValue: true}
	command.Action = Run
	parent.Commands = append(parent.Commands, command)
}
