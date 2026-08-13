package main

import (
	spotdepth "github.com/k4k3ru-hub/binance/go/cmd/cli/rest/spot/depth"
	spotexchangeinfo "github.com/k4k3ru-hub/binance/go/cmd/cli/rest/spot/exchange_info"
	usdsmdepth "github.com/k4k3ru-hub/binance/go/cmd/cli/rest/usdsm/depth"
	usdsmexchangeinfo "github.com/k4k3ru-hub/binance/go/cmd/cli/rest/usdsm/exchange_info"
	spotwsdepth "github.com/k4k3ru-hub/binance/go/cmd/cli/ws/spot/depth"
	usdsmwsdepth "github.com/k4k3ru-hub/binance/go/cmd/cli/ws/usdsm/depth"

	"github.com/k4k3ru-hub/cli-go"
)

const (
	restCommandName  = "rest"
	spotCommandName  = "spot"
	usdsmCommandName = "usdsm"
	wsCommandName    = "ws"
)

func main() {
	app := cli.NewCli(nil)
	app.SetVersion("0.1.0")
	app.Command.SetDefaultConfigOption()

	restCommand := cli.NewCommand(restCommandName)
	restCommand.Usage = "REST API commands."
	app.Command.Commands = append(app.Command.Commands, restCommand)

	spotCommand := cli.NewCommand(spotCommandName)
	spotCommand.Usage = "Binance Spot commands."
	restCommand.Commands = append(restCommand.Commands, spotCommand)

	usdsmCommand := cli.NewCommand(usdsmCommandName)
	usdsmCommand.Usage = "Binance USDⓈ-M Futures commands."
	restCommand.Commands = append(restCommand.Commands, usdsmCommand)

	spotexchangeinfo.SetCommand(spotCommand)
	spotdepth.SetCommand(spotCommand)
	usdsmexchangeinfo.SetCommand(usdsmCommand)
	usdsmdepth.SetCommand(usdsmCommand)

	wsCommand := cli.NewCommand(wsCommandName)
	wsCommand.Usage = "WebSocket stream commands."
	app.Command.Commands = append(app.Command.Commands, wsCommand)

	wsSpotCommand := cli.NewCommand(spotCommandName)
	wsSpotCommand.Usage = "Binance Spot streams."
	wsCommand.Commands = append(wsCommand.Commands, wsSpotCommand)
	spotwsdepth.SetCommand(wsSpotCommand)

	wsUSDSMCommand := cli.NewCommand(usdsmCommandName)
	wsUSDSMCommand.Usage = "Binance USDⓈ-M Futures streams."
	wsCommand.Commands = append(wsCommand.Commands, wsUSDSMCommand)
	usdsmwsdepth.SetCommand(wsUSDSMCommand)

	app.Run()
}
