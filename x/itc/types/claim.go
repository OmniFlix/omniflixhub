package types

import (
	"github.com/OmniFlix/omniflixhub/x/itc/exported"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
)

var (
	_ proto.Message   = &Claim{}
	_ exported.ClaimI = &Claim{}
)

func NewClaim(campaignId uint64, address string) Claim {
	return Claim{
		CampaignId: campaignId,
		Address:    address,
	}
}

func (c Claim) GetCampaignId() uint64 {
	return c.CampaignId
}

func (c Claim) GetAddress() sdk.AccAddress {
	address, _ := sdk.AccAddressFromBech32(c.Address)
	return address
}
