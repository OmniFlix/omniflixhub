package main

import (
	"os"

	"github.com/OmniFlix/omniflixhub/app"
	"github.com/OmniFlix/omniflixhub/cmd/omniflixhubd/cmd"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, "OMNIFLIXHUBD", app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
