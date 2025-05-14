package cli

import (
	"context"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/OmniFlix/omniflixhub/v6/x/onft/types"
)

func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                types.ModuleName,
		Short:              "Querying commands for the oNFT module",
		DisableFlagParsing: true,
	}

	queryCmd.AddCommand(
		GetCmdQueryDenom(),
		GetCmdQueryDenoms(),
		GetCmdQueryCollection(),
		GetCmdQuerySupply(),
		GetCmdQueryONFT(),
		GetCmdQueryOwner(),
		GetCmdQueryParams(),
	)

	return queryCmd
}

func GetCmdQuerySupply() *cobra.Command {
	cmd := &cobra.Command{
		Use: "supply [denom-id]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`total supply of a collection of oNFTs.
Example:
$ %s query onft supply [denom-id]`, version.AppName)),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			var owner sdk.AccAddress
			ownerStr, err := cmd.Flags().GetString(FlagOwner)
			if err != nil {
				return err
			}

			if len(ownerStr) > 0 {
				owner, err = sdk.AccAddressFromBech32(ownerStr)
				if err != nil {
					return err
				}
			}

			denomId := args[0]

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Supply(context.Background(), &types.QuerySupplyRequest{
				DenomId: denomId,
				Owner:   owner.String(),
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}
	cmd.Flags().AddFlagSet(FsQuerySupply)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "owner [address]",
		Long:    "Get the oNFTs owned by an account address.",
		Example: fmt.Sprintf("$ %s query onft owner <addr> --denom-id=<denom-id>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if _, err := sdk.AccAddressFromBech32(args[0]); err != nil {
				return err
			}
			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			denomID, err := cmd.Flags().GetString(FlagDenomID)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.OwnerONFTs(context.Background(), &types.QueryOwnerONFTsRequest{
				DenomId:    denomID,
				Owner:      args[0],
				Pagination: pagination,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}
	cmd.Flags().AddFlagSet(FsQueryOwner)
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "owner onfts")

	return cmd
}

func GetCmdQueryCollection() *cobra.Command {
	cmd := &cobra.Command{
		Use: "collection [denom-id]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Get all the oNFTs from a given collection
Example:
$ %s query onft collection <denom-id>`, version.AppName)),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			denomId := args[0]
			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Collection(context.Background(), &types.QueryCollectionRequest{
				DenomId:    denomId,
				Pagination: pagination,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp.Collection)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "onfts")
	return cmd
}

func GetCmdQueryDenoms() *cobra.Command {
	cmd := &cobra.Command{
		Use: "denoms",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all denominations of all collections of oNFTs
Example:
$ %s query onft denoms`, version.AppName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			owner, err := cmd.Flags().GetString(FlagOwner)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Denoms(
				context.Background(),
				&types.QueryDenomsRequest{
					Pagination: pagination,
					Owner:      owner,
				},
			)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "denoms")
	cmd.Flags().String(FlagOwner, "", "filter by collection owner address")
	return cmd
}

func GetCmdQueryDenom() *cobra.Command {
	cmd := &cobra.Command{
		Use: "denom [denom-id]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the denominations by the specified denom name
Example:
$ %s query onft denom <denom-id>`, version.AppName)),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			denomId := args[0]

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Denom(context.Background(), &types.QueryDenomRequest{
				DenomId: denomId,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp.Denom)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryONFT() *cobra.Command {
	cmd := &cobra.Command{
		Use: "asset [denom-id] [onft-id]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query a single oNFT from a collection
Example:
$ %s query onft asset <denom> <onft-id>`, version.AppName)),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			denomId := args[0]
			onftId := args[1]

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.ONFT(context.Background(), &types.QueryONFTRequest{
				DenomId: denomId,
				Id:      onftId,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp.ONFT)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query oNFT params",
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
