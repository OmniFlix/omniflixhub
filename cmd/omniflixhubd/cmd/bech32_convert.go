package cmd

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/bech32"

	"github.com/spf13/cobra"
)

var flagBech32Prefix = "prefix"

// AddBech32ConvertCommand returns bech32-convert cobra Command.
func AddBech32ConvertCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bech32-convert [address]",
		Short: "Convert any bech32 string to the omniflix prefix",
		Long: `Convert any bech32 string to the omniflix prefix

Example:
	omniflixhubd debug bech32-convert juno1a6zlyvpnksx8wr6wz8wemur2xe8zyh0yst53ep

	omniflixhubd debug bech32-convert cosmos1673f0t8p893rqyqe420mgwwz92ac4qv6n0nsjx --prefix omniflix
	`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			bech32prefix, err := cmd.Flags().GetString(flagBech32Prefix)
			if err != nil {
				return err
			}

			address := args[0]
			convertedAddress, err := ConvertBech32Prefix(address, bech32prefix)
			if err != nil {
				return fmt.Errorf("convertation failed: %s", err)
			}

			cmd.Println(convertedAddress)

			return nil
		},
	}

	cmd.Flags().StringP(flagBech32Prefix, "p", "cosmos", "Bech32 Prefix to encode to")

	return cmd
}

func ConvertBech32Prefix(address, prefix string) (string, error) {
	_, bz, err := bech32.DecodeAndConvert(address)
	if err != nil {
		return "", fmt.Errorf("cannot decode %s address: %s", address, err)
	}

	convertedAddress, err := bech32.ConvertAndEncode(prefix, bz)
	if err != nil {
		return "", fmt.Errorf("cannot convert %s address: %s", address, err)
	}

	return convertedAddress, nil
}

// addDebugCommands adds the custom debug commands to the application.
func addDebugCommands(cmd *cobra.Command) *cobra.Command {
	cmd.AddCommand(AddBech32ConvertCommand())
	return cmd
}
