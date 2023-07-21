package main

import (
	"os"

	"github.com/OmniFlix/omniflixhub/v2/app"
	"github.com/OmniFlix/omniflixhub/v2/cmd/omniflixhubd/cmd"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
