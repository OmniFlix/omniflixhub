package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/OmniFlix/omniflixhub/v6/x/medianode/types"
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
		GetCmdQueryMediaNodes(),
	)

	return cmd
}

// GetCmdQueryCampaign implements the query campaign command.
func GetCmdQueryMediaNode() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "medianode [id]",
		Long:    "Query a campaign by it's id.",
		Example: fmt.Sprintf("$ %s query medianode medianode <id>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			mediaNodeId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.MediaNode(context.Background(), &types.QueryMediaNodeRequest{
				Id: mediaNodeId,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.MediaNode)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryAllCampaigns implements the query all campaigns command.
func GetCmdQueryMediaNodes() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "medianodes",
		Long:    "Query medianodes.",
		Example: fmt.Sprintf("$ %s query medianode medianodes", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			resp, err := queryClient.MediaNodes(
				context.Background(),
				&types.QueryMediaNodesRequest{
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
	flags.AddPaginationFlagsToCmd(cmd, "media-nodes")

	return cmd
}
