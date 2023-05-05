package cli

import (
	"fmt"
	"strconv"

	"github.com/OmniFlix/omniflixhub/x/itc/types"
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
		GetCmdCreateCampaign(),
		GetCmdCancelCampaign(),
		GetCmdDepositCampaign(),
		GetCmdClaim(),
	)

	return itcTxCmd
}

// GetCmdCreateCampaign implements the create-campaign command
func GetCmdCreateCampaign() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-campaign",
		Short: "creates a campaign on itc module",
		Example: fmt.Sprintf(
			"$ %s tx itc create-campaign "+
				"--name=<name> "+
				"--description=<description> "+
				"--start-time=<start-time> "+
				"--duration=<duration> "+
				"--claim-type=<claim-type> "+
				"--max-allowed-claims=<max-claims> "+
				"--tokens-per-claim=<tokens-per-claim> "+
				"--deposit=<token deposit> "+
				"--distribution-type=<distr-type> "+
				"--interaction-type=<interaction-type> "+
				"--nft-denom-id=<denom-id> "+
				"--nft-details-file=<path/to/nft-details> "+
				"--from=<key-name> "+
				"--chain-id=<chain-id> "+
				"--fees=<fee> ",
			version.AppName,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			creator := clientCtx.GetFromAddress()
			name, description, err := parseCampaignNameAndDescription(cmd.Flags())
			if err != nil {
				return err
			}
			claimType, err := parseClaimType(cmd.Flags())
			if err != nil {
				return err
			}
			interactionType, err := parseInteractionType(cmd.Flags())
			if err != nil {
				return err
			}
			nftDenomId, err := cmd.Flags().GetString(FlagNftDenomId)
			if err != nil {
				return err
			}
			maxAllowedClaims, err := cmd.Flags().GetUint64(FlagMaxAllowedClaims)
			if err != nil {
				return err
			}
			startTime, duration, err := parseStartTimeAndDuration(cmd.Flags())
			if err != nil {
				return err
			}
			tokensPerClaim, tokensDeposited, err := parseCampaignTokens(cmd.Flags(), claimType)
			if err != nil {
				return err
			}
			distribution, err := parseDistribution(cmd.Flags())
			if err != nil {
				return err
			}
			nftDetails, err := parseNftDetails(cmd.Flags(), claimType)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateCampaign(name,
				description,
				interactionType,
				claimType,
				nftDenomId,
				maxAllowedClaims,
				tokensPerClaim,
				tokensDeposited,
				nftDetails,
				distribution,
				startTime,
				duration,
				creator.String(),
			)
			fmt.Printf("%v", msg)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsCreateCampaign)
	_ = cmd.MarkFlagRequired(FlagNftDenomId)
	_ = cmd.MarkFlagRequired(FlagName)
	_ = cmd.MarkFlagRequired(FlagDescription)
	_ = cmd.MarkFlagRequired(FlagClaimType)
	_ = cmd.MarkFlagRequired(FlagInteractionType)
	_ = cmd.MarkFlagRequired(FlagStartTime)
	_ = cmd.MarkFlagRequired(FlagDuration)
	_ = cmd.MarkFlagRequired(FlagMaxAllowedClaims)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdCancelCampaign implements the campaign cancel command
func GetCmdCancelCampaign() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-campaign",
		Short: "cancels the campaign before start",
		Example: fmt.Sprintf(
			"$ %s tx itc cancel-campaign [campaign-id] "+
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

			creator := clientCtx.GetFromAddress()
			campaignId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgCancelCampaign(campaignId, creator.String())

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdDepositCampaign implements the bid command
func GetCmdDepositCampaign() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "campaign-deposit",
		Short: "deposits tokens into a campaign",
		Example: fmt.Sprintf(
			"$ %s tx itc deposit [campaign-id] "+
				"--amount=<amount> "+
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

			depositor := clientCtx.GetFromAddress()
			campaignId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			amountStr, err := cmd.Flags().GetString(FlagAmount)
			if err != nil {
				return err
			}
			amount, err := sdk.ParseCoinNormalized(amountStr)
			if err != nil {
				return err
			}

			msg := types.NewMsgDepositCampaign(campaignId, amount, depositor.String())

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsDepositCampaign)
	_ = cmd.MarkFlagRequired(FlagAmount)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdClaim implements the bid command
func GetCmdClaim() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim",
		Short: "claim tokens from a campaign",
		Example: fmt.Sprintf(
			"$ %s tx itc claim [campaign-id] "+
				"--nft-id=<nft-id> "+
				"--interaction-type=<interaction> "+
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

			claimer := clientCtx.GetFromAddress()
			campaignId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			nftId, err := cmd.Flags().GetString(FlagNftId)
			if err != nil {
				return err
			}
			interactType, err := parseInteractionType(cmd.Flags())

			msg := types.NewMsgClaim(campaignId, nftId, interactType, claimer.String())

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsClaim)
	_ = cmd.MarkFlagRequired(FlagNftId)
	_ = cmd.MarkFlagRequired(FlagInteractionType)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
