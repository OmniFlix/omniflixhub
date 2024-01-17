package v2

import (
	"github.com/OmniFlix/omniflixhub/v2/app/keepers"
	"github.com/OmniFlix/omniflixhub/v2/app/upgrades"
	itctypes "github.com/OmniFlix/omniflixhub/v2/x/itc/types"
	marketplacetypes "github.com/OmniFlix/omniflixhub/v2/x/marketplace/types"
	onfttypes "github.com/OmniFlix/omniflixhub/v2/x/onft/types"
	streampaytypes "github.com/OmniFlix/streampay/v2/x/streampay/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	icqtypes "github.com/cosmos/ibc-apps/modules/async-icq/v7/types"
	icahosttypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/types"
)

func CreateV2UpgradeHandler(
	mm *module.Manager,
	cfg module.Configurator,
	_ upgrades.BaseAppParamManager,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("running migrations ...")

		// Run migrations before applying any other state changes.
		// NOTE: DO NOT PUT ANY STATE CHANGES BEFORE RunMigrations().
		versionMap, err := mm.RunMigrations(ctx, cfg, fromVM)
		if err != nil {
			return nil, err
		}
		// set streampay default params
		err = keepers.StreamPayKeeper.SetParams(ctx, streampaytypes.DefaultParams())
		if err != nil {
			return nil, err
		}

		hostParams := icahosttypes.Params{
			HostEnabled: true,
			AllowMessages: []string{
				sdk.MsgTypeURL(&banktypes.MsgSend{}),
				sdk.MsgTypeURL(&stakingtypes.MsgDelegate{}),
				sdk.MsgTypeURL(&stakingtypes.MsgUndelegate{}),
				sdk.MsgTypeURL(&stakingtypes.MsgBeginRedelegate{}),
				sdk.MsgTypeURL(&stakingtypes.MsgCreateValidator{}),
				sdk.MsgTypeURL(&stakingtypes.MsgEditValidator{}),
				sdk.MsgTypeURL(&distrtypes.MsgWithdrawDelegatorReward{}),
				sdk.MsgTypeURL(&distrtypes.MsgWithdrawValidatorCommission{}),
				sdk.MsgTypeURL(&distrtypes.MsgSetWithdrawAddress{}),
				sdk.MsgTypeURL(&distrtypes.MsgFundCommunityPool{}),
				sdk.MsgTypeURL(&govv1beta1.MsgVote{}),
				// onft
				sdk.MsgTypeURL(&onfttypes.MsgCreateDenom{}),
				sdk.MsgTypeURL(&onfttypes.MsgUpdateDenom{}),
				sdk.MsgTypeURL(&onfttypes.MsgTransferDenom{}),
				sdk.MsgTypeURL(&onfttypes.MsgMintONFT{}),
				sdk.MsgTypeURL(&onfttypes.MsgTransferONFT{}),
				sdk.MsgTypeURL(&onfttypes.MsgBurnONFT{}),
				// marketplace
				sdk.MsgTypeURL(&marketplacetypes.MsgListNFT{}),
				sdk.MsgTypeURL(&marketplacetypes.MsgDeListNFT{}),
				sdk.MsgTypeURL(&marketplacetypes.MsgBuyNFT{}),
				sdk.MsgTypeURL(&marketplacetypes.MsgCreateAuction{}),
				sdk.MsgTypeURL(&marketplacetypes.MsgCancelAuction{}),
				sdk.MsgTypeURL(&marketplacetypes.MsgPlaceBid{}),
				// itc
				sdk.MsgTypeURL(&itctypes.MsgCreateCampaign{}),
				sdk.MsgTypeURL(&itctypes.MsgCancelCampaign{}),
				sdk.MsgTypeURL(&itctypes.MsgDepositCampaign{}),
				sdk.MsgTypeURL(&itctypes.MsgClaim{}),
				// streampay
				sdk.MsgTypeURL(&streampaytypes.MsgStreamSend{}),
				sdk.MsgTypeURL(&streampaytypes.MsgStopStream{}),
				sdk.MsgTypeURL(&streampaytypes.MsgClaimStreamedAmount{}),
			},
		}
		keepers.ICAHostKeeper.SetParams(ctx, hostParams)

		// set icq queries
		icqQueries := icqtypes.Params{
			HostEnabled: true,
			AllowQueries: []string{
				"/ibc.applications.transfer.v1.Query/DenomTrace",
				"/cosmos.auth.v1beta1.Query/Account",
				"/cosmos.auth.v1beta1.Query/Params",
				"/cosmos.bank.v1beta1.Query/Balance",
				"/cosmos.bank.v1beta1.Query/DenomMetadata",
				"/cosmos.bank.v1beta1.Query/Params",
				"/cosmos.bank.v1beta1.Query/SupplyOf",
				"/cosmos.distribution.v1beta1.Query/Params",
				"/cosmos.distribution.v1beta1.Query/DelegatorWithdrawAddress",
				"/cosmos.distribution.v1beta1.Query/ValidatorCommission",
				"/cosmos.gov.v1beta1.Query/Deposit",
				"/cosmos.gov.v1beta1.Query/Params",
				"/cosmos.gov.v1beta1.Query/Vote",
				"/cosmos.slashing.v1beta1.Query/Params",
				"/cosmos.slashing.v1beta1.Query/SigningInfo",
				"/cosmos.staking.v1beta1.Query/Delegation",
				"/cosmos.staking.v1beta1.Query/Params",
				"/cosmos.staking.v1beta1.Query/Validator",
				"/osmosis.tokenfactory.v1beta1.Query/Params",
				"/osmosis.tokenfactory.v1beta1.Query/DenomAuthorityMetadata",
				"/ibc.applications.nft_transfer.v1.Query/Params",
				"/ibc.applications.nft_transfer.v1.Query/ClassTrace",
				"/ibc.applications.nft_transfer.v1.Query/ClassHash",
				"/ibc.applications.nft_transfer.v1.Query/EscrowAddress",

				"/OmniFlix.onft.v1beta1.Query/Params",
				"/OmniFlix.onft.v1beta1.Query/Denom",
				"/OmniFlix.onft.v1beta1.Query/ONFT",
				"/OmniFlix.onft.v1beta1.Query/Supply",

				"/OmniFlix.marketplace.v1beta1.Query/Params",
				"/OmniFlix.marketplace.v1beta1.Query/Listing",
				"/OmniFlix.marketplace.v1beta1.Query/Auction",

				"/OmniFlix.itc.v1.Query/Params",
				"/OmniFlix.itc.v1.Query/Campaign",

				"/OmniFlix.streampay.v1.Query/Params",
				"/OmniFlix.streampay.v1.Query/StreamingPayment",
			},
		}
		_ = keepers.ICQKeeper.SetParams(ctx, icqQueries)

		ctx.Logger().Info("Upgrade complete")
		return versionMap, nil
	}
}
