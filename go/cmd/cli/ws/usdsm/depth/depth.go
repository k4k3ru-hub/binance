package depth

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/k4k3ru-hub/binance/go/websocket"
	wsprotocol "github.com/k4k3ru-hub/binance/go/websocket/protocol"
	usdsmdepth "github.com/k4k3ru-hub/binance/go/websocket/usdsm/depth"
	"github.com/k4k3ru-hub/cli-go"
)

const (
	CommandName            = "depth"
	OptionNameSymbol       = "symbol"
	OptionAliasSymbol      = "s"
	OptionNameUpdateSpeed  = "update-speed"
	OptionAliasUpdateSpeed = "u"
)

var command = cli.NewCommand(CommandName)

func ParamsFromOptions(options map[string]*cli.Option) (usdsmdepth.Params, error) {
	params := usdsmdepth.Params{}
	if option := options[OptionNameSymbol]; option != nil {
		params.Symbol = strings.ToUpper(strings.TrimSpace(option.Value))
	}
	if params.Symbol == "" {
		return params, fmt.Errorf("symbol is required")
	}
	if option := options[OptionNameUpdateSpeed]; option != nil {
		params.UpdateSpeed = usdsmdepth.UpdateSpeed(strings.ToLower(strings.TrimSpace(option.Value)))
	}
	if _, err := params.Stream(); err != nil {
		return params, err
	}
	return params, nil
}

type sessionHandler struct{}

func (*sessionHandler) HandleMessage(_ websocket.SessionContext, message []byte) {
	var event wsprotocol.USDSMDepthEvent
	if err := json.Unmarshal(message, &event); err != nil || event.EventType != "depthUpdate" {
		return
	}
	if err := json.NewEncoder(os.Stdout).Encode(event); err != nil {
		fmt.Printf("failed to output USDⓈ-M depth event: %v\n", err)
	}
}

func (*sessionHandler) HandleClose(websocket.SessionContext) {}

func Run(options map[string]*cli.Option) {
	params, err := ParamsFromOptions(options)
	if err != nil {
		fmt.Printf("%s\n", err)
		command.ShowUsage()
		return
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	client, err := websocket.NewUSDSMClient(ctx, &sessionHandler{}, nil)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	defer client.Close()
	if err := client.Depth().Subscribe(ctx, params); err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	<-ctx.Done()
}

func SetCommand(parent *cli.Command) {
	command.Usage = "Subscribe to USDⓈ-M Futures diff-depth updates."
	command.Options[OptionNameSymbol] = &cli.Option{Alias: OptionAliasSymbol, HasValue: true}
	command.Options[OptionNameUpdateSpeed] = &cli.Option{Alias: OptionAliasUpdateSpeed, HasValue: true}
	command.Action = Run
	parent.Commands = append(parent.Commands, command)
}
