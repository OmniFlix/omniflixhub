package main

import (
	"os"

	"github.com/OmniFlix/omniflixhub/app"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/tendermint/starport/starport/pkg/cosmoscmd"
	tmcmds "github.com/tendermint/tendermint/cmd/tendermint/commands"
)

/*
func main() {
	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
*/

func main() {
	cmdOptions := GetWasmCmdOptions()
	cmdOptions = append(cmdOptions, cosmoscmd.AddSubCmd(tmcmds.RollbackStateCmd))
	rootCmd, _ := cosmoscmd.NewRootCmd(
		app.Name,
		app.AccountAddressPrefix,
		app.DefaultNodeHome,
		app.Name,
		app.ModuleBasics,
		app.NewOmniFlixApp,
		// this line is used by starport scaffolding # root/arguments
		cmdOptions...,
	)
	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
