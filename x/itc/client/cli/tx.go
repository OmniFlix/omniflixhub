package cli

import (
	"fmt"
	"strconv"
	"strings"
	"time"

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
		GetCmdCampaignDeposit(),
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
				"--claimable-tokens=<claimable-tokens> "+
				"--total-tokens=<total-tokens> "+
				"--distribution-type=<distr-type> "+
				"--interaction-type=<interaction-type> "+
				"--nft-denom-id=<denom-id> "+
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
			name, err := cmd.Flags().GetString(FlagName)
			if err != nil {
				return err
			}
			description, err := cmd.Flags().GetString(FlagDescription)
			if err != nil {
				return err
			}
			claimType, err := cmd.Flags().GetString(FlagClaimType)
			if err != nil {
				return err
			}
			var claim types.ClaimType
			if claimType == "fungible" {
				claim = types.CLAIM_TYPE_FT
			} else if claimType == "non-fungible" {
				claim = types.CLAIM_TYPE_NFT
			} else if claimType == "fungible-and-non-fungible" {
				claim = types.CLAIM_TYPE_FT_AND_NFT
			}
			interactType, err := cmd.Flags().GetString(FlagInteractionType)
			if err != nil {
				return err
			}
			interactType = strings.ToLower(interactType)
			var interaction types.InteractionType
			if interactType == "hold" {
				interaction = types.INTERACTION_TYPE_HOLD
			} else if interactType == "burn" {
				interaction = types.INTERACTION_TYPE_BURN
			} else if interactType == "transfer" {
				interaction = types.INTERACTION_TYPE_TRANSFER
			}
			nftDenomId, err := cmd.Flags().GetString(FlagNftDenomId)
			if err != nil {
				return err
			}
			maxAllowedClaims, err := cmd.Flags().GetUint64(FlagMaxAllowedClaims)
			if err != nil {
				return err
			}
			claimableTokensStr, err := cmd.Flags().GetString(FlagClaimableTokens)
			if err != nil {
				return err
			}
			claimTokens, err := sdk.ParseCoinNormalized(claimableTokensStr)
			if err != nil {
				return err
			}
			claimableTokens := types.Tokens{
				Fungible: &claimTokens,
			}
			totalTokensStr, err := cmd.Flags().GetString(FlagTotalTokens)
			if err != nil {
				return err
			}
			tokensDeposit, err := sdk.ParseCoinNormalized(totalTokensStr)
			if err != nil {
				return err
			}
			totalTokens := types.Tokens{
				Fungible: &tokensDeposit,
			}
			_, err = cmd.Flags().GetString(FlagDistributionType)
			if err != nil {
				return err
			}
			distribution := &types.Distribution{
				Type: types.DISTRIBUTION_TYPE_INSTANT,
			}
			startTimeStr, err := cmd.Flags().GetString(FlagStartTime)
			if err != nil {
				return err
			}
			var startTime time.Time
			if startTimeStr != "" {
				startTime, err = time.Parse(time.RFC3339, startTimeStr)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("failed to parse start time: %s", startTime)
			}
			durationStr, err := cmd.Flags().GetString(FlagDuration)
			if err != nil {
				return err
			}
			duration, err := time.ParseDuration(durationStr)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateCampaign(name,
				description,
				interaction,
				claim,
				nftDenomId,
				maxAllowedClaims,
				claimableTokens,
				totalTokens,
				nil,
				distribution,
				startTime,
				duration,
				creator.String(),
			)

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
	_ = cmd.MarkFlagRequired(FlagClaimableTokens)
	_ = cmd.MarkFlagRequired(FlagTotalTokens)
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

// GetCmdCampaignDeposit implements the bid command
func GetCmdCampaignDeposit() *cobra.Command {
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

			msg := types.NewMsgCampaignDeposit(campaignId, amount, depositor.String())

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsCampaignDeposit)
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
			interaction, err := cmd.Flags().GetString(FlagInteractionType)
			if err != nil {
				return err
			}
			interaction = strings.ToLower(interaction)
			var interactType types.InteractionType
			if interaction == "hold" {
				interactType = types.INTERACTION_TYPE_HOLD
			} else if interaction == "transfer" {
				interactType = types.INTERACTION_TYPE_TRANSFER
			} else if interaction == "burn" {
				interactType = types.INTERACTION_TYPE_BURN
			}

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
