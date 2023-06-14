package keeper

import (
	"context"
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/OmniFlix/omniflixhub/x/itc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.QueryServer = Keeper{}

// Params returns params of the itc module.
func (k Keeper) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	var params types.Params
	k.paramSpace.GetParamSet(ctx, &params)

	return &types.QueryParamsResponse{Params: params}, nil
}

// Campaign returns details of the campaign
func (k Keeper) Campaign(goCtx context.Context,
	req *types.QueryCampaignRequest,
) (*types.QueryCampaignResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	campaign, found := k.GetCampaign(ctx, req.CampaignId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "campaign %d not found", req.CampaignId)
	}

	return &types.QueryCampaignResponse{Campaign: campaign}, nil
}

func (k Keeper) Campaigns(goCtx context.Context,
	req *types.QueryCampaignsRequest,
) (*types.QueryCampaignsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	var filteredCampaigns []types.Campaign
	var pageRes *query.PageResponse
	store := ctx.KVStore(k.storeKey)
	campaignStore := prefix.NewStore(store, types.PrefixCampaignId)
	pageRes, err := query.FilteredPaginate(campaignStore,
		req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
			var al types.Campaign
			k.cdc.MustUnmarshal(value, &al)
			matchCreator, matchStatus := true, true
			// match status (if supplied/valid)
			if len(req.Status.String()) > 0 && types.ValidCampaignStatus(req.Status) {
				if req.Status == types.CAMPAIGN_STATUS_ACTIVE {
					matchStatus = al.StartTime.Before(time.Now())
				} else if req.Status == types.CAMPAIGN_STATUS_INACTIVE {
					matchStatus = al.StartTime.After(time.Now())
				}
			}

			// match creator address (if supplied)
			if len(req.Creator) > 0 {
				creator, err := sdk.AccAddressFromBech32(req.Creator)
				if err != nil {
					return false, err
				}

				matchCreator = al.Creator == creator.String()
			}

			if matchCreator && matchStatus {
				if accumulate {
					filteredCampaigns = append(filteredCampaigns, al)
				}

				return true, nil
			}

			return false, nil
		})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
	}

	return &types.QueryCampaignsResponse{Campaigns: filteredCampaigns, Pagination: pageRes}, nil
}

func (k Keeper) Claims(goCtx context.Context,
	req *types.QueryClaimsRequest,
) (*types.QueryClaimsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	var filteredClaims []types.Claim
	var pageRes *query.PageResponse
	store := ctx.KVStore(k.storeKey)
	claimStore := prefix.NewStore(store, types.PrefixClaimByNftId)
	pageRes, err := query.FilteredPaginate(claimStore,
		req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
			var al types.Claim
			k.cdc.MustUnmarshal(value, &al)
			matchClaimer := true
			matchCampaignId := al.CampaignId == req.CampaignId
			// match claimer address (if supplied)
			if len(req.Address) > 0 {
				address, err := sdk.AccAddressFromBech32(req.Address)
				if err != nil {
					return false, err
				}

				matchClaimer = al.Address == address.String()
			}

			if matchCampaignId && matchClaimer {
				if accumulate {
					filteredClaims = append(filteredClaims, al)
				}

				return true, nil
			}

			return false, nil
		})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
	}
	return &types.QueryClaimsResponse{
		Claims:     filteredClaims,
		Pagination: pageRes,
	}, nil
}
