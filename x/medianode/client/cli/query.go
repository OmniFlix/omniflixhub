package cli

import (
	"context"
	"fmt"

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
		GetCmdQueryParams(),
		GetCmdQueryMediaNodes(),
		GetCmdQueryMediaNode(),
		GetCmdQueryLease(),
		GetCmdQueryLeasesByLessee(),
	)

	return cmd
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query medianode params",
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

func GetCmdQueryMediaNode() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "node [id]",
		Long:    "Query a medianode by it's id.",
		Example: fmt.Sprintf("$ %s query medianode node <id>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			mediaNodeId := args[0]

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

func GetCmdQueryMediaNodes() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "nodes",
		Long:    "Query medianodes.",
		Example: fmt.Sprintf("$ %s query medianode nodes", version.AppName),
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
	flags.AddPaginationFlagsToCmd(cmd, "nodes")

	return cmd
}

// GetCmdQueryLease implements the query lease command.
func GetCmdQueryLease() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "lease [id]",
		Long:    "Query a lease by its id.",
		Example: fmt.Sprintf("$ %s query medianode lease <id>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			mediaNodeId := args[0]

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Lease(context.Background(), &types.QueryLeaseRequest{
				MediaNodeId: mediaNodeId,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Lease)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryLeasesByLessee implements the query all leases command.
func GetCmdQueryLeasesByLessee() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "leases-by-lessee",
		Long:    "Query leases by lessee.",
		Example: fmt.Sprintf("$ %s query medianode leases-by-lessee", version.AppName),
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
			resp, err := queryClient.LeasesByLessee(
				context.Background(),
				&types.QueryLeasesByLesseeRequest{
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
	flags.AddPaginationFlagsToCmd(cmd, "leases-by-lessee")

	return cmd
}
