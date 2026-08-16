package depth

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/k4k3ru-hub/binance/go/rest"
	restprotocol "github.com/k4k3ru-hub/binance/go/rest/protocol"
	spotdepth "github.com/k4k3ru-hub/binance/go/rest/spot/depth"
	"github.com/k4k3ru-hub/cli-go"
)

const (
	CommandName            = "depth"
	OptionNameSymbol       = "symbol"
	OptionAliasSymbol      = "s"
	OptionNameLimit        = "limit"
	OptionAliasLimit       = "l"
	OptionNameSymbolStatus = "symbol-status"
)

var command = cli.NewCommand(CommandName)

// ParamsFromOptions validates CLI options and returns immutable API parameters.
func ParamsFromOptions(options map[string]*cli.Option) (spotdepth.Params, error) {
	params := spotdepth.Params{}
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
	if option := options[OptionNameSymbolStatus]; option != nil {
		params.SymbolStatus = strings.ToUpper(strings.TrimSpace(option.Value))
	}
	return params, nil
}

// Run executes the Spot depth command.
func Run(options map[string]*cli.Option) {
	params, err := ParamsFromOptions(options)
	if err != nil {
		fmt.Printf("%s\n", err)
		command.ShowUsage()
		return
	}
	client, err := rest.NewClient(nil)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	result, err := client.Spot().Depth().Send(context.Background(), params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	outputBook(result.Bids, result.Asks)
}

func outputBook(bids, asks []restprotocol.PriceLevel) {
	headers := []string{"Side", "Price", "Quantity"}
	data := make([][]interface{}, 0, len(bids)+len(asks))
	for _, level := range bids {
		data = append(data, []interface{}{"BID", level.Price, level.Quantity})
	}
	for _, level := range asks {
		data = append(data, []interface{}{"ASK", level.Price, level.Quantity})
	}
	cli.OutputTable(headers, data)
}

// SetCommand registers the command with its parent.
func SetCommand(parent *cli.Command) {
	command.Usage = "Get a Spot order-book snapshot."
	command.Options[OptionNameSymbol] = &cli.Option{Alias: OptionAliasSymbol, HasValue: true}
	command.Options[OptionNameLimit] = &cli.Option{Alias: OptionAliasLimit, HasValue: true}
	command.Options[OptionNameSymbolStatus] = &cli.Option{HasValue: true}
	command.Action = Run
	parent.Commands = append(parent.Commands, command)
}
