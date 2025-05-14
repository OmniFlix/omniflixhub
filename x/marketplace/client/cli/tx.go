package cli

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/OmniFlix/omniflixhub/v6/x/marketplace/types"
	"github.com/cosmos/cosmos-sdk/client"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	marketplaceTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	marketplaceTxCmd.AddCommand(
		GetCmdListNft(),
		GetCmdEditListing(),
		GetCmdDeListNft(),
		GetCmdBuyNft(),
		GetCmdCreateAuction(),
		GetCmdCancelAuction(),
		GetCmdPlaceBid(),
	)

	return marketplaceTxCmd
}

// GetCmdListNft implements the list-nft command
func GetCmdListNft() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "list-nft",
		Long: "lists an nft on marketplace",
		Example: fmt.Sprintf(
			"$ %s tx marketplace list-nft "+
				"--nft-id=<nft-id> "+
				"--denom-id=<nft-id> "+
				"--price=\"1000000uflix\" "+
				"--from=<key-name> "+
				"--chain-id=<chain-id> "+
				"--fees=<fee>",
			version.AppName,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			owner := clientCtx.GetFromAddress()
			denomId, err := cmd.Flags().GetString(FlagDenomId)
			if err != nil {
				return err
			}
			nftId, err := cmd.Flags().GetString(FlagNftId)
			if err != nil {
				return err
			}
			priceStr, err := cmd.Flags().GetString(FlagPrice)
			if err != nil {
				return err
			}
			price, err := sdk.ParseCoinNormalized(priceStr)
			if err != nil {
				return fmt.Errorf("failed to parse price: %s", price)
			}
			splitSharesStr, err := cmd.Flags().GetString(FlagSplitShares)
			if err != nil {
				return err
			}
			var splitShares []types.WeightedAddress
			if len(splitSharesStr) > 0 {
				splitShares, err = parseSplitShares(splitSharesStr)
				if err != nil {
					return err
				}
			}

			msg := types.NewMsgListNFT(denomId, nftId, price, owner, splitShares)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsListNft)
	_ = cmd.MarkFlagRequired(FlagDenomId)
	_ = cmd.MarkFlagRequired(FlagNftId)
	_ = cmd.MarkFlagRequired(FlagPrice)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdEditListing implements the edit-listing command
func GetCmdEditListing() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "edit-listing",
		Long: "Edit an existing marketplace listing ",
		Example: fmt.Sprintf(
			"$ %s tx marketplace edit-listing [listing-id] "+
				"--price=\"1000000uflix\" "+
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

			owner := clientCtx.GetFromAddress()

			listingId := args[0]

			priceStr, err := cmd.Flags().GetString(FlagPrice)
			if err != nil {
				return err
			}
			price, err := sdk.ParseCoinNormalized(priceStr)
			if err != nil {
				return fmt.Errorf("failed to parse price: %s", price)
			}

			msg := types.NewMsgEditListing(listingId, price, owner)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsEditListing)
	_ = cmd.MarkFlagRequired(FlagPrice)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdDeListNft implements the de-list-nft command
func GetCmdDeListNft() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "de-list-nft",
		Long: "de-list an existing listing from marketplace",
		Example: fmt.Sprintf(
			"$ %s tx marketplace de-list-nft [listing-id] "+
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

			owner := clientCtx.GetFromAddress()

			listingId := args[0]

			msg := types.NewMsgDeListNFT(listingId, owner)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdBuyNft implements the buy-nft command
func GetCmdBuyNft() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "buy-nft",
		Short: "Buy an nft from marketplace",
		Example: fmt.Sprintf(
			"$ %s tx marketplace buy-nft [listing-id]"+
				"--price=<price>"+
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

			buyer := clientCtx.GetFromAddress()
			listingId := args[0]

			priceStr, err := cmd.Flags().GetString(FlagPrice)
			if err != nil {
				return err
			}
			price, err := sdk.ParseCoinNormalized(priceStr)
			if err != nil {
				return fmt.Errorf("failed to parse price: %s", priceStr)
			}

			msg := types.NewMsgBuyNFT(listingId, price, buyer)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsBuyNFT)
	_ = cmd.MarkFlagRequired(FlagPrice)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func parseSplitShares(splitsharesStr string) ([]types.WeightedAddress, error) {
	splitsharesStr = strings.TrimSpace(splitsharesStr)
	splitsStrList := strings.Split(splitsharesStr, ",")
	var weightedAddrsList []types.WeightedAddress
	for _, splitStr := range splitsStrList {
		var share types.WeightedAddress
		split := strings.Split(strings.TrimSpace(splitStr), ":")
		address, err := sdk.AccAddressFromBech32(strings.TrimSpace(split[0]))
		if err != nil {
			return nil, err
		}
		weight, err := sdkmath.LegacyNewDecFromStr(strings.TrimSpace(split[1]))
		if err != nil {
			return nil, err
		}
		share.Address = address.String()
		share.Weight = weight
		weightedAddrsList = append(weightedAddrsList, share)
	}
	return weightedAddrsList, nil
}

func parseWhitelistAccounts(whitelistStr string) ([]string, error) {
	whitelistStr = strings.TrimSpace(whitelistStr)
	whitelist := strings.Split(whitelistStr, ",")
	return whitelist, nil
}

// GetCmdCreateAuction implements the create-auction command
func GetCmdCreateAuction() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create-auction",
		Long: "creates an auction on marketplace",
		Example: fmt.Sprintf(
			"$ %s tx marketplace create-auction "+
				"--nft-id=<nft-id> "+
				"--denom-id=<nft-id> "+
				"--start-price=\"1000000uflix\" "+
				"--start-time=\"2022-06-13T13:02:49.389Z\" "+
				"--increment-percentage=\"0.01\""+
				"--from=<key-name> "+
				"--chain-id=<chain-id> "+
				"--fees=<fee>",
			version.AppName,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			owner := clientCtx.GetFromAddress()
			denomId, err := cmd.Flags().GetString(FlagDenomId)
			if err != nil {
				return err
			}
			nftId, err := cmd.Flags().GetString(FlagNftId)
			if err != nil {
				return err
			}
			startPriceStr, err := cmd.Flags().GetString(FlagStartPrice)
			if err != nil {
				return err
			}
			startPrice, err := sdk.ParseCoinNormalized(startPriceStr)
			if err != nil {
				return fmt.Errorf("failed to parse start price: %s", startPrice)
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
			splitSharesStr, err := cmd.Flags().GetString(FlagSplitShares)
			if err != nil {
				return err
			}
			var splitShares []types.WeightedAddress
			if len(splitSharesStr) > 0 {
				splitShares, err = parseSplitShares(splitSharesStr)
				if err != nil {
					return err
				}
			}
			durationStr, err := cmd.Flags().GetString(FlagDuration)
			if err != nil {
				return err
			}
			var duration *time.Duration
			if len(durationStr) > 0 {
				dur, err := time.ParseDuration(durationStr)
				if err != nil {
					return err
				}
				duration = &dur
			} else {
				duration = nil
			}
			incrementStr, err := cmd.Flags().GetString(FlagIncrementPercentage)
			if err != nil {
				return err
			}
			increment, err := sdkmath.LegacyNewDecFromStr(incrementStr)
			if err != nil {
				return err
			}
			whitelistAccountsStr, err := cmd.Flags().GetString(FlagWhiteListAccounts)
			if err != nil {
				return err
			}
			var whitelist []string
			if len(whitelistAccountsStr) > 0 {
				whitelist, err = parseWhitelistAccounts(whitelistAccountsStr)
				if err != nil {
					return err
				}
			}
			msg := types.NewMsgCreateAuction(denomId, nftId, startTime, duration, startPrice, owner, increment, whitelist, splitShares)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsCreateAuction)
	_ = cmd.MarkFlagRequired(FlagDenomId)
	_ = cmd.MarkFlagRequired(FlagNftId)
	_ = cmd.MarkFlagRequired(FlagStartPrice)
	_ = cmd.MarkFlagRequired(FlagStartTime)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdCancelAuction implements the cancel-auction command
func GetCmdCancelAuction() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "cancel-auction",
		Long: "cancel an existing auction from marketplace with no bids",
		Example: fmt.Sprintf(
			"$ %s tx marketplace cancel-auction [auction-id] "+
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

			owner := clientCtx.GetFromAddress()

			auctionId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgCancelAuction(auctionId, owner)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdBid implements the bid command
func GetCmdPlaceBid() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "place-bid",
		Short: "Bid for an nft on marketplace",
		Example: fmt.Sprintf(
			"$ %s tx marketplace place-bid [auction-id]"+
				"--amount=<amount>"+
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

			buyer := clientCtx.GetFromAddress()
			auctionId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			amountStr, err := cmd.Flags().GetString(FlagAmount)
			if err != nil {
				return err
			}
			amount, err := sdk.ParseCoinNormalized(amountStr)
			if err != nil {
				return fmt.Errorf("failed to parse price: %s", amountStr)
			}

			msg := types.NewMsgPlaceBid(auctionId, amount, buyer)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsPlaceBid)
	_ = cmd.MarkFlagRequired(FlagAmount)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
