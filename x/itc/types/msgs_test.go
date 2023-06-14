package types_test

import (
	"testing"
	"time"

	"github.com/OmniFlix/omniflixhub/x/itc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

var (
	defaultCampaignName        = "test campaign"
	defaultCampaignDescription = "test description"
	defaultNftDenomId          = "onftdenomtest001"
	defaultNftId               = "onfttest001"
	defaultCampaignId          = uint64(1)
	defaultTokensPerClaim      = sdk.NewInt64Coin(types.DefaultCampaignCreationFee.Denom, 10_000_000)
	defaultDuration            = time.Second * 100
	defaultDistribution        = types.Distribution{Type: types.DISTRIBUTION_TYPE_STREAM, StreamDuration: time.Second * 300}
)

func TestNewMsgCreateCampaign(t *testing.T) {
	// generate a private/public key pair and get the respective address
	pk1 := ed25519.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pk1.Address())

	// make a valid create campaign message
	createCampaignMsg := func(
		after func(msg types.MsgCreateCampaign) types.MsgCreateCampaign,
		interaction types.InteractionType,
		claimType types.ClaimType,
		nftMintDetails *types.NFTDetails,
		distribution *types.Distribution,
	) types.MsgCreateCampaign {
		validCreateCampaignMsg := *types.NewMsgCreateCampaign(
			defaultCampaignName,
			defaultCampaignDescription,
			interaction,
			claimType,
			defaultNftDenomId,
			10,
			sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 10_000_000),
			sdk.NewInt64Coin(defaultTokensPerClaim.Denom, 100_000_000),
			nftMintDetails,
			distribution,
			time.Now(),
			defaultDuration,
			addr1.String(),
			types.DefaultCampaignCreationFee,
		)

		return after(validCreateCampaignMsg)
	}
	msg := createCampaignMsg(func(msg types.MsgCreateCampaign) types.MsgCreateCampaign {
		return msg
	}, types.INTERACTION_TYPE_BURN, types.CLAIM_TYPE_FT, nil, &defaultDistribution)

	require.Equal(t, msg.Route(), types.RouterKey)
	require.Equal(t, msg.Type(), "create_campaign")
	signers := msg.GetSigners()
	require.Equal(t, len(signers), 1)
	require.Equal(t, signers[0].String(), addr1.String())

	testCases := []struct {
		desc  string
		msg   types.MsgCreateCampaign
		valid bool
	}{
		{
			desc: "valid message",
			msg: createCampaignMsg(func(msg types.MsgCreateCampaign) types.MsgCreateCampaign {
				return msg
			}, types.INTERACTION_TYPE_BURN, types.CLAIM_TYPE_FT, nil, &defaultDistribution),
			valid: true,
		},
		{
			desc: "invalid creator",
			msg: createCampaignMsg(func(msg types.MsgCreateCampaign) types.MsgCreateCampaign {
				msg.Creator = ""
				return msg
			}, types.INTERACTION_TYPE_BURN, types.CLAIM_TYPE_FT, nil, &defaultDistribution),
			valid: false,
		},
		{
			desc: "invalid interaction type",
			msg: createCampaignMsg(func(msg types.MsgCreateCampaign) types.MsgCreateCampaign {
				return msg
			}, 5, types.CLAIM_TYPE_FT, nil, &defaultDistribution),
			valid: false,
		},
		{
			desc: "invalid claim type",
			msg: createCampaignMsg(func(msg types.MsgCreateCampaign) types.MsgCreateCampaign {
				return msg
			}, types.INTERACTION_TYPE_BURN, 5, nil, &defaultDistribution),
			valid: false,
		},
		{
			desc: "invalid duration",
			msg: createCampaignMsg(func(msg types.MsgCreateCampaign) types.MsgCreateCampaign {
				msg.Duration = time.Second * 0
				return msg
			}, types.INTERACTION_TYPE_BURN, types.CLAIM_TYPE_FT, nil, &defaultDistribution),
			valid: false,
		},
		{
			desc: "invalid deposit amount",
			msg: createCampaignMsg(func(msg types.MsgCreateCampaign) types.MsgCreateCampaign {
				msg.Deposit = sdk.NewInt64Coin("test", 0)
				return msg
			}, types.INTERACTION_TYPE_BURN, types.CLAIM_TYPE_FT, nil, &defaultDistribution),
			valid: false,
		},
		{
			desc: "nil deposit amount",
			msg: createCampaignMsg(func(msg types.MsgCreateCampaign) types.MsgCreateCampaign {
				msg.Deposit = sdk.Coin{}
				return msg
			}, types.INTERACTION_TYPE_BURN, types.CLAIM_TYPE_FT, nil, &defaultDistribution),
			valid: false,
		},
		{
			desc: "invalid tokens per claim",
			msg: createCampaignMsg(func(msg types.MsgCreateCampaign) types.MsgCreateCampaign {
				msg.TokensPerClaim = sdk.NewInt64Coin("test", 0)
				return msg
			}, types.INTERACTION_TYPE_BURN, types.CLAIM_TYPE_FT, nil, &defaultDistribution),
			valid: false,
		},
		{
			desc: "invalid distribution",
			msg: createCampaignMsg(func(msg types.MsgCreateCampaign) types.MsgCreateCampaign {
				msg.Distribution = &types.Distribution{
					Type: types.DISTRIBUTION_TYPE_STREAM,
				}
				return msg
			}, types.INTERACTION_TYPE_BURN, types.CLAIM_TYPE_FT, nil, &defaultDistribution),
			valid: false,
		},
		{
			desc: "invalid nft mint details",
			msg: createCampaignMsg(func(msg types.MsgCreateCampaign) types.MsgCreateCampaign {
				return msg
			}, types.INTERACTION_TYPE_BURN, types.CLAIM_TYPE_NFT, nil, nil),
		},
		{
			desc: "invalid nft mint details - empty media uri",
			msg: createCampaignMsg(func(msg types.MsgCreateCampaign) types.MsgCreateCampaign {
				msg.NftMintDetails = &types.NFTDetails{
					DenomId:    defaultNftDenomId,
					Name:       "test",
					PreviewUri: "ipfs://test",
				}
				return msg
			}, types.INTERACTION_TYPE_BURN, types.CLAIM_TYPE_NFT, nil, nil),
			valid: false,
		},
	}
	for _, tc := range testCases {
		if tc.valid {
			require.NoError(t, tc.msg.ValidateBasic(), "test: %v", tc.desc)
		} else {
			require.Error(t, tc.msg.ValidateBasic(), "test: %v", tc.desc)
		}
	}
}

func TestNewMsgCancelCampaign(t *testing.T) {
	// generate a private/public key pair and get the respective address
	pk1 := ed25519.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pk1.Address())

	cancelCampaignMsg := types.NewMsgCancelCampaign(defaultCampaignId, addr1.String())
	require.Equal(t, cancelCampaignMsg.Route(), types.RouterKey)
	require.Equal(t, cancelCampaignMsg.Type(), "cancel_campaign")
	signers := cancelCampaignMsg.GetSigners()
	require.Equal(t, len(signers), 1)
	require.Equal(t, signers[0].String(), addr1.String())

	testCases := []struct {
		desc  string
		msg   types.MsgCancelCampaign
		valid bool
	}{
		{
			desc:  "valid message",
			msg:   *types.NewMsgCancelCampaign(defaultCampaignId, addr1.String()),
			valid: true,
		},
		{
			desc:  "invalid creator",
			msg:   *types.NewMsgCancelCampaign(defaultCampaignId, ""),
			valid: false,
		},
	}

	for _, tc := range testCases {
		if tc.valid {
			require.NoError(t, tc.msg.ValidateBasic(), "test: %v", tc.desc)
		} else {
			require.Error(t, tc.msg.ValidateBasic(), "test: %v", tc.desc)
		}
	}
}

func TestNewMsgClaim(t *testing.T) {
	// generate a private/public key pair and get the respective address
	pk1 := ed25519.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pk1.Address())

	claimMsg := types.NewMsgClaim(defaultCampaignId, defaultNftId, types.INTERACTION_TYPE_BURN, addr1.String())
	require.Equal(t, claimMsg.Route(), types.RouterKey)
	require.Equal(t, claimMsg.Type(), "claim")
	signers := claimMsg.GetSigners()
	require.Equal(t, len(signers), 1)
	require.Equal(t, signers[0].String(), addr1.String())

	testCases := []struct {
		desc  string
		msg   types.MsgClaim
		valid bool
	}{
		{
			desc:  "valid message",
			msg:   *types.NewMsgClaim(defaultCampaignId, addr1.String(), types.INTERACTION_TYPE_BURN, addr1.String()),
			valid: true,
		},
		{
			desc:  "invalid claimer",
			msg:   *types.NewMsgClaim(defaultCampaignId, defaultNftId, types.INTERACTION_TYPE_BURN, "fioenfor"),
			valid: false,
		},
		{
			desc:  "invalid nft id",
			msg:   *types.NewMsgClaim(defaultCampaignId, "", types.INTERACTION_TYPE_BURN, addr1.String()),
			valid: false,
		},
		{
			desc:  "invalid interaction type",
			msg:   *types.NewMsgClaim(defaultCampaignId, defaultNftId, 5, addr1.String()),
			valid: false,
		},
	}

	for _, tc := range testCases {
		if tc.valid {
			require.NoError(t, tc.msg.ValidateBasic(), "test: %v", tc.desc)
		} else {
			require.Error(t, tc.msg.ValidateBasic(), "test: %v", tc.desc)
		}
	}
}

func TestNewMsgDepositCampaign(t *testing.T) {
	// generate a private/public key pair and get the respective address
	pk1 := ed25519.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pk1.Address())

	depositCampaignMsg := types.NewMsgDepositCampaign(defaultCampaignId, defaultTokensPerClaim, addr1.String())
	require.Equal(t, depositCampaignMsg.Route(), types.RouterKey)
	require.Equal(t, depositCampaignMsg.Type(), "deposit_campaign")
	signers := depositCampaignMsg.GetSigners()
	require.Equal(t, len(signers), 1)
	require.Equal(t, signers[0].String(), addr1.String())

	testCases := []struct {
		desc  string
		msg   types.MsgDepositCampaign
		valid bool
	}{
		{
			desc:  "valid message",
			msg:   *types.NewMsgDepositCampaign(defaultCampaignId, defaultTokensPerClaim, addr1.String()),
			valid: true,
		},
		{
			desc: "invalid depositor",
			msg: *types.NewMsgDepositCampaign(
				defaultCampaignId,
				sdk.NewInt64Coin("test", 100),
				"fwngoew",
			),
			valid: false,
		},
		{
			desc: "invalid deposit",
			msg: *types.NewMsgDepositCampaign(
				defaultCampaignId,
				sdk.Coin{},
				addr1.String(),
			),
			valid: false,
		},
	}

	for _, tc := range testCases {
		if tc.valid {
			require.NoError(t, tc.msg.ValidateBasic(), "test: %v", tc.desc)
		} else {
			require.Error(t, tc.msg.ValidateBasic(), "test: %v", tc.desc)
		}
	}
}
