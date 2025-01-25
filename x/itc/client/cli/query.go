package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/OmniFlix/omniflixhub/v6/x/itc/types"
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
		GetCmdQueryCampaign(),
		GetCmdQueryAllCampaigns(),
		GetCmdQueryClaimsByCampaign(),
	)

	return cmd
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query itc params",
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

// GetCmdQueryCampaign implements the query campaign command.
func GetCmdQueryCampaign() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "campaign [id]",
		Long:    "Query a campaign by it's id.",
		Example: fmt.Sprintf("$ %s query itc campaign <id>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			campaignId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Campaign(context.Background(), &types.QueryCampaignRequest{
				CampaignId: campaignId,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Campaign)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryAllCampaigns implements the query all campaigns command.
func GetCmdQueryAllCampaigns() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "campaigns",
		Long:    "Query campaigns.",
		Example: fmt.Sprintf("$ %s query itc campaigns", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			creator, err := cmd.Flags().GetString(FlagCreator)
			if err != nil {
				return err
			}
			campaignStatus, err := parseCampaignStatus(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			resp, err := queryClient.Campaigns(
				context.Background(),
				&types.QueryCampaignsRequest{
					Creator:    creator,
					Status:     campaignStatus,
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
	cmd.Flags().String(FlagCreator, "", "filter by creator address")
	cmd.Flags().String(FlagStatus, "", "filter by campaign status")
	flags.AddPaginationFlagsToCmd(cmd, "all campaigns")

	return cmd
}

// GetCmdQueryClaimsByCampaign implements the query claims by campaign command.
func GetCmdQueryClaimsByCampaign() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "claims [campaign-id]",
		Long:    "Query claims by campaign.",
		Example: fmt.Sprintf("$ %s query itc claims <campaign-id>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			claimerStr, err := cmd.Flags().GetString(FlagClaimer)
			if err != nil {
				return err
			}
			var claimer sdk.AccAddress
			if len(claimerStr) > 0 {
				claimer, err = sdk.AccAddressFromBech32(claimerStr)
				if err != nil {
					return err
				}
			}
			campaignId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			resp, err := queryClient.Claims(
				context.Background(),
				&types.QueryClaimsRequest{
					Address:    claimer.String(),
					CampaignId: campaignId,
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
	cmd.Flags().String(FlagClaimer, "", "filter by claimer address")
	flags.AddPaginationFlagsToCmd(cmd, "campaign claims")

	return cmd
}
