package types

import (
	"github.com/OmniFlix/omniflixhub/v4/x/itc/exported"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
)

var (
	_ proto.Message   = &Claim{}
	_ exported.ClaimI = &Claim{}
)

func NewClaim(campaignId uint64, address, nftId string, interactionType InteractionType) Claim {
	return Claim{
		CampaignId:  campaignId,
		Address:     address,
		NftId:       nftId,
		Interaction: interactionType,
	}
}

func (c Claim) GetCampaignId() uint64 {
	return c.CampaignId
}

func (c Claim) GetNftId() string {
	return c.NftId
}

func (c Claim) GetInteractionType() string {
	return c.Interaction.String()
}

func (c Claim) GetAddress() sdk.AccAddress {
	address, err := sdk.AccAddressFromBech32(c.Address)
	if err != nil {
		return nil
	}
	return address
}
