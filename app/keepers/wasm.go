package keepers

import (
	"strings"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	itctypes "github.com/OmniFlix/omniflixhub/v4/x/itc/types"
	marketplacetypes "github.com/OmniFlix/omniflixhub/v4/x/marketplace/types"
	onfttypes "github.com/OmniFlix/omniflixhub/v4/x/onft/types"
	tokenfactorytypes "github.com/OmniFlix/omniflixhub/v4/x/tokenfactory/types"
	streampaytypes "github.com/OmniFlix/streampay/v2/x/streampay/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
)

// AllCapabilities returns all capabilities available with the current wasmvm
// See https://github.com/CosmWasm/cosmwasm/blob/main/docs/CAPABILITIES-BUILT-IN.md
// This functionality is going to be moved upstream: https://github.com/CosmWasm/wasmvm/issues/425
var wasmCapabilities = []string{
	"iterator",
	"staking",
	"stargate",
	"cosmwasm_1_1",
	"cosmwasm_1_2",
	"cosmwasm_1_3",
	"cosmwasm_1_4",
	"cosmwasm_1_5",
	"token_factory",
}

func AcceptedStargateQueries() wasmkeeper.AcceptedStargateQueries {
	return wasmkeeper.AcceptedStargateQueries{
		// ibc
		"/ibc.core.client.v1.Query/ClientState":    &ibcclienttypes.QueryClientStateResponse{},
		"/ibc.core.client.v1.Query/ConsensusState": &ibcclienttypes.QueryConsensusStateResponse{},
		"/ibc.core.connection.v1.Query/Connection": &ibcconnectiontypes.QueryConnectionResponse{},

		// governance
		"/cosmos.gov.v1beta1.Query/Vote": &govv1.QueryVoteResponse{},

		// distribution
		"/cosmos.distribution.v1beta1.Query/DelegationRewards": &distrtypes.QueryDelegationRewardsResponse{},

		// staking
		"/cosmos.staking.v1beta1.Query/Delegation":          &stakingtypes.QueryDelegationResponse{},
		"/cosmos.staking.v1beta1.Query/Redelegations":       &stakingtypes.QueryRedelegationsResponse{},
		"/cosmos.staking.v1beta1.Query/UnbondingDelegation": &stakingtypes.QueryUnbondingDelegationResponse{},
		"/cosmos.staking.v1beta1.Query/Validator":           &stakingtypes.QueryValidatorResponse{},
		"/cosmos.staking.v1beta1.Query/Params":              &stakingtypes.QueryParamsResponse{},
		"/cosmos.staking.v1beta1.Query/Pool":                &stakingtypes.QueryPoolResponse{},

		// onft
		"/OmniFlix.onft.v1beta1.Query/Denoms":        &onfttypes.QueryDenomsResponse{},
		"/OmniFlix.onft.v1beta1.Query/Denom":         &onfttypes.QueryDenomResponse{},
		"/OmniFlix.onft.v1beta1.Query/IBCDenom":      &onfttypes.QueryDenomResponse{},
		"/OmniFlix.onft.v1beta1.Query/Collection":    &onfttypes.QueryCollectionResponse{},
		"/OmniFlix.onft.v1beta1.Query/IBCCollection": &onfttypes.QueryCollectionResponse{},
		"/OmniFlix.onft.v1beta1.Query/OwnerONFTs":    &onfttypes.QueryOwnerONFTsResponse{},
		"/OmniFlix.onft.v1beta1.Query/ONFT":          &onfttypes.QueryONFTResponse{},
		"/OmniFlix.onft.v1beta1.Query/Supply":        &onfttypes.QuerySupplyResponse{},
		"/OmniFlix.onft.v1beta1.Query/Params":        &onfttypes.QueryParamsResponse{},

		// marketplace
		"/OmniFlix.marketplace.v1beta1.Query/Listings":        &marketplacetypes.QueryListingsResponse{},
		"/OmniFlix.marketplace.v1beta1.Query/Listing":         &marketplacetypes.QueryListingResponse{},
		"/OmniFlix.marketplace.v1beta1.Query/ListingsByOwner": &marketplacetypes.QueryListingsByOwnerResponse{},
		"/OmniFlix.marketplace.v1beta1.Query/Auctions":        &marketplacetypes.QueryAuctionsResponse{},
		"/OmniFlix.marketplace.v1beta1.Query/Auction":         &marketplacetypes.QueryAuctionResponse{},
		"/OmniFlix.marketplace.v1beta1.Query/AuctionsByOwner": &marketplacetypes.QueryAuctionsResponse{},
		"/OmniFlix.marketplace.v1beta1.Query/Bids":            &marketplacetypes.QueryBidsResponse{},
		"/OmniFlix.marketplace.v1beta1.Query/Bid":             &marketplacetypes.QueryBidResponse{},
		"/OmniFlix.marketplace.v1beta1.Query/Params":          &marketplacetypes.QueryParamsResponse{},

		// itc
		"/OmniFlix.itc.v1.Query/Campaigns": &itctypes.QueryCampaignsResponse{},
		"/OmniFlix.itc.v1.Query/Campaign":  &itctypes.QueryCampaignResponse{},
		"/OmniFlix.itc.v1.Query/Claims":    &itctypes.QueryClaimsResponse{},
		"/OmniFlix.itc.v1.Query/Params":    &itctypes.QueryParamsResponse{},

		// streampay
		"/OmniFlix.streampay.v1.Query/StreamPayments": &streampaytypes.QueryStreamPaymentsResponse{},
		"/OmniFlix.streampay.v1.Query/StreamPayment":  &streampaytypes.QueryStreamPaymentResponse{},
		"/OmniFlix.streampay.v1.Query/Params":         &streampaytypes.QueryParamsResponse{},

		// tokenfactory queries
		"/osmosis.tokenfactory.v1beta1.Query/Params":                 &tokenfactorytypes.QueryParamsResponse{},
		"/osmosis.tokenfactory.v1beta1.Query/DenomAuthorityMetadata": &tokenfactorytypes.QueryDenomAuthorityMetadataResponse{},
		"/osmosis.tokenfactory.v1beta1.Query/DenomsFromCreator":      &tokenfactorytypes.QueryDenomsFromCreatorResponse{},
	}
}

func GetWasmCapabilities() string {
	return strings.Join(wasmCapabilities, ",")
}
