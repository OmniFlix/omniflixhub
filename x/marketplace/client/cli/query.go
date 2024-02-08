package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/OmniFlix/omniflixhub/v3/x/marketplace/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	// Group marketplace queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(
		GetCmdQueryParams(),
		GetCmdQueryListing(),
		GetCmdQueryAllListings(),
		GetCmdQueryListingsByOwner(),
		GetCmdQueryAuction(),
		GetCmdQueryAllAuctions(),
		GetCmdQueryAuctionsByOwner(),
		GetCmdQueryAuctionBid(),
		GetCmdQueryAllBids(),
	)

	return cmd
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query Marketplace params",
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryListing implements the query listing command.
func GetCmdQueryListing() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "listing [id]",
		Long:    "Query a listing by id.",
		Example: fmt.Sprintf("$ %s query marketplace listing <id>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			listingId := args[0]

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Listing(context.Background(), &types.QueryListingRequest{
				Id: listingId,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res.Listing)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryAllListings implements the query all listings command.
func GetCmdQueryAllListings() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "listings",
		Long:    "Query listings.",
		Example: fmt.Sprintf("$ %s query marketplace listings", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			owner, err := cmd.Flags().GetString(FlagOwner)
			if err != nil {
				return err
			}
			priceDenom, err := cmd.Flags().GetString(FlagPriceDenom)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			resp, err := queryClient.Listings(
				context.Background(),
				&types.QueryListingsRequest{
					Owner:      owner,
					PriceDenom: priceDenom,
					Pagination: pageReq,
				},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	cmd.Flags().String(FlagOwner, "", "filter by owner address")
	cmd.Flags().String(FlagPriceDenom, "", "filter by listing price-denom")
	flags.AddPaginationFlagsToCmd(cmd, "all listings")

	return cmd
}

// GetCmdQueryListingsByOwner implements the query listings by owner command.
func GetCmdQueryListingsByOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "listings-by-owner [owner]",
		Long:    "Query listings by the owner.",
		Example: fmt.Sprintf("$ %s query marketplace listings <owner>", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			var owner sdk.AccAddress
			if len(args) > 0 {
				owner, err = sdk.AccAddressFromBech32(args[0])
				if err != nil {
					return err
				}
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			resp, err := queryClient.ListingsByOwner(
				context.Background(),
				&types.QueryListingsByOwnerRequest{
					Owner:      owner.String(),
					Pagination: pageReq,
				},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "owner listings")

	return cmd
}

// GetCmdQueryAuctionListing implements the query auction command.
func GetCmdQueryAuction() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "auction [id]",
		Long:    "Query a auction by id.",
		Example: fmt.Sprintf("$ %s query marketplace auction <id>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			auctionId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Auction(context.Background(), &types.QueryAuctionRequest{
				Id: auctionId,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res.Auction)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryAllAuctions implements the query all auctions command.
func GetCmdQueryAllAuctions() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "auctions",
		Long:    "Query auctions.",
		Example: fmt.Sprintf("$ %s query marketplace auctions", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			owner, err := cmd.Flags().GetString(FlagOwner)
			if err != nil {
				return err
			}
			priceDenom, err := cmd.Flags().GetString(FlagPriceDenom)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			resp, err := queryClient.Auctions(
				context.Background(),
				&types.QueryAuctionsRequest{
					Owner:      owner,
					PriceDenom: priceDenom,
					Pagination: pageReq,
				},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	cmd.Flags().String(FlagOwner, "", "filter by owner address")
	cmd.Flags().String(FlagPriceDenom, "", "filter by auction price-denom")
	flags.AddPaginationFlagsToCmd(cmd, "all auctions")

	return cmd
}

// GetCmdQueryAuctionsByOwner implements the query auctions by owner command.
func GetCmdQueryAuctionsByOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "auctions-by-owner [owner]",
		Long:    "Query auctions by the owner.",
		Example: fmt.Sprintf("$ %s query marketplace auctions <owner>", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			var owner sdk.AccAddress
			if len(args) > 0 {
				owner, err = sdk.AccAddressFromBech32(args[0])
				if err != nil {
					return err
				}
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			resp, err := queryClient.AuctionsByOwner(
				context.Background(),
				&types.QueryAuctionsByOwnerRequest{
					Owner:      owner.String(),
					Pagination: pageReq,
				},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "owner auctions")

	return cmd
}

// GetCmdQueryAuctionBid implements the query bid command.
func GetCmdQueryAuctionBid() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bid [id]",
		Long:    "Query a bid by auction id.",
		Example: fmt.Sprintf("$ %s query marketplace bid <id>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			auctionId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Bid(context.Background(), &types.QueryBidRequest{
				Id: auctionId,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res.Bid)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryAllBids implements the query all bids command.
func GetCmdQueryAllBids() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bids",
		Long:    "Query bids.",
		Example: fmt.Sprintf("$ %s query marketplace bids", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			bidder, err := cmd.Flags().GetString(FlagBidder)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			resp, err := queryClient.Bids(
				context.Background(),
				&types.QueryBidsRequest{
					Bidder:     bidder,
					Pagination: pageReq,
				},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	cmd.Flags().String(FlagBidder, "", "filter by bidder address")
	flags.AddPaginationFlagsToCmd(cmd, "all bids")

	return cmd
}
