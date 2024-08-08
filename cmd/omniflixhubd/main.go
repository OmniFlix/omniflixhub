package main

import (
	"os"

	"github.com/OmniFlix/omniflixhub/v5/app"
	"github.com/OmniFlix/omniflixhub/v5/cmd/omniflixhubd/cmd"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
)

func main() {
	app.SetConfig()

	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
