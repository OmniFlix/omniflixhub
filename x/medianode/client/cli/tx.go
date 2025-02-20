package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/OmniFlix/omniflixhub/v6/x/medianode/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	itcTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	itcTxCmd.AddCommand(
		GetCmdRegisterMediaNode(),
		GetCmdLeaseMediaNode(),
		GetCmdDepositMediaNode(),
		GetCmdCancelLease(),
		GetCmdCloseMediaNode(),
		GetCmdExtendLease(),
	)

	return itcTxCmd
}

// GetCmdRegisterMediaNode implements the register-media-node command
func GetCmdRegisterMediaNode() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register",
		Short: "registers a new media node",
		Long:  "register a new media node with the specified URL, hardware specifications, and lease price per day\n",
		Example: fmt.Sprintf(
			"$ %s tx medianode register [url] --hardware-specs=<hardwarespecs> --price-per-day=<amount> "+
				"--from=<key-name> "+
				"--chain-id=<chain-id> "+
				"--fees=<fee>",
			version.AppName,
		),
		Args: cobra.ExactArgs(1), // Expecting 1 positional argument
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			url := args[0]
			hardwareSpecsStr, err := cmd.Flags().GetString(FlagHardwareSpecs)
			if err != nil {
				return err
			}
			hardwareSpecs, err := ParseHardwareSpecs(hardwareSpecsStr)
			if err != nil {
				return err
			}
			priceStr, err := cmd.Flags().GetString(FlagPricePerHour)
			if err != nil {
				return err
			}
			price, err := sdk.ParseCoinNormalized(priceStr)
			if err != nil {
				return err
			}
			depositAmountStr, err := cmd.Flags().GetString(FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinNormalized(depositAmountStr)
			if err != nil {
				return err
			}

			msg, err := types.NewMsgRegisterMediaNode(url, hardwareSpecs, price, deposit, clientCtx.GetFromAddress().String())
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsRegisterMediaNode) // Assuming FsRegisterMediaNode is defined elsewhere
	_ = cmd.MarkFlagRequired(FlagHardwareSpecs)
	_ = cmd.MarkFlagRequired(FlagPricePerHour)
	_ = cmd.MarkFlagRequired(FlagDeposit)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdLeaseMediaNode implements the lease-media-node command
func GetCmdLeaseMediaNode() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lease",
		Short: "leases a media node",
		Long:  "leases a media node with the specified URL and lease hours\n",
		Example: fmt.Sprintf(
			"$ %s tx medianode lease [medianode-id] --lease-hours=<no-of-hours> --amount=<amount> "+
				"--from=<key-name> "+
				"--chain-id=<chain-id> "+
				"--fees=<fee>",
			version.AppName,
		),
		Args: cobra.ExactArgs(1), // Expecting 1 positional argument
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			mediaNodeId := args[0]
			leaseDays, err := cmd.Flags().GetUint64(FlagLeaseHours)
			if err != nil {
				return err
			}
			leaseAmountStr, err := cmd.Flags().GetString(FlagLeaseAmount)
			if err != nil {
				return err
			}
			leaseAmount, err := sdk.ParseCoinNormalized(leaseAmountStr)
			if err != nil {
				return err
			}

			msg := types.NewMsgLeaseMediaNode(mediaNodeId, leaseDays, leaseAmount, clientCtx.GetFromAddress().String())

			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsLeaseMediaNode)
	_ = cmd.MarkFlagRequired(FlagLeaseHours)
	_ = cmd.MarkFlagRequired(FlagLeaseAmount)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdExtendLease implements the extend-lease command
func GetCmdExtendLease() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extend-lease",
		Short: "extends a lease for a media node",
		Long:  "extends an active lease for a media node with the specified ID and additional lease hours\n",
		Example: fmt.Sprintf(
			"$ %s tx medianode extend-lease [medianode-id] --additional-hours=<no-of-hours> --amount=<amount> "+
				"--from=<key-name> "+
				"--chain-id=<chain-id> "+
				"--fees=<fee>",
			version.AppName,
		),
		Args: cobra.ExactArgs(1), // Expecting 1 positional argument
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			mediaNodeId := args[0]
			additionalHours, err := cmd.Flags().GetUint64(FlagLeaseHours)
			if err != nil {
				return err
			}
			leaseAmountStr, err := cmd.Flags().GetString(FlagLeaseAmount)
			if err != nil {
				return err
			}
			leaseAmount, err := sdk.ParseCoinNormalized(leaseAmountStr)
			if err != nil {
				return err
			}

			msg := types.NewMsgExtendLease(mediaNodeId, additionalHours, leaseAmount, clientCtx.GetFromAddress().String())

			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsExtendLease)
	_ = cmd.MarkFlagRequired(FlagLeaseHours)
	_ = cmd.MarkFlagRequired(FlagLeaseAmount)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdDepositMediaNode implements the deposit-media-node command
func GetCmdDepositMediaNode() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit",
		Short: "deposits an amount for a media node",
		Long:  "deposits an amount for a media node with the specified URL and deposit amount\n",
		Example: fmt.Sprintf(
			"$ %s tx medianode deposit [medianode-id] --amount=<amount> "+
				"--from=<key-name> "+
				"--chain-id=<chain-id> "+
				"--fees=<fee>",
			version.AppName,
		),
		Args: cobra.ExactArgs(1), // Expecting 1 positional argument
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			medianodeId := args[0]

			depositStr, err := cmd.Flags().GetString(FlagDeposit)
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinNormalized(depositStr)
			if err != nil {
				return err
			}

			msg := types.NewMsgDepositMediaNode(medianodeId, deposit, clientCtx.GetFromAddress().String())

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsDepositMediaNode)
	_ = cmd.MarkFlagRequired(FlagDeposit)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdCancelLease implements the cancel-lease command
func GetCmdCancelLease() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-lease",
		Short: "cancels a lease for a media node",
		Long:  "cancels an active lease for a media node with the specified ID\n",
		Example: fmt.Sprintf(
			"$ %s tx medianode cancel-lease [medianode-id] "+
				"--from=<key-name> "+
				"--chain-id=<chain-id> "+
				"--fees=<fee>",
			version.AppName,
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			mediaNodeId := args[0]

			msg := types.NewMsgCancelLease(mediaNodeId, clientCtx.GetFromAddress().String())

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdCloseMediaNode implements the close-media-node command
func GetCmdCloseMediaNode() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "close",
		Short: "closes a media node",
		Long:  "closes a media node with the specified ID\n",
		Example: fmt.Sprintf(
			"$ %s tx medianode close [medianode-id] "+
				"--from=<key-name> "+
				"--chain-id=<chain-id> "+
				"--fees=<fee>",
			version.AppName,
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			mediaNodeId := args[0]

			msg := types.NewMsgCloseMediaNode(mediaNodeId, clientCtx.GetFromAddress().String())

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// ParseHardwareSpecs extracts hardware specifications from a comma-separated string
func ParseHardwareSpecs(hardwareSpecsStr string) (types.HardwareSpecs, error) {
	specsStr := strings.Split(hardwareSpecsStr, ",")
	if len(specsStr) != 3 {
		return types.HardwareSpecs{}, fmt.Errorf("expected 3 hardware specs: cpus, ram, storage")
	}

	var specs types.HardwareSpecs

	// Parse CPUs
	cpus, err := strconv.ParseInt(strings.TrimSpace(specsStr[0]), 10, 64)
	if err != nil {
		return specs, fmt.Errorf("invalid CPU spec: %s", specsStr[0])
	}
	specs.Cpus = cpus

	// Parse RAM
	ram, err := strconv.ParseInt(strings.TrimSpace(specsStr[1]), 10, 64)
	if err != nil {
		return specs, fmt.Errorf("invalid RAM spec: %s", specsStr[1])
	}
	specs.RamInGb = ram

	// Parse Storage
	storage, err := strconv.ParseInt(strings.TrimSpace(specsStr[2]), 10, 64)
	if err != nil {
		return specs, fmt.Errorf("invalid Storage spec: %s", specsStr[2])
	}
	specs.StorageInGb = storage

	return specs, nil
}
