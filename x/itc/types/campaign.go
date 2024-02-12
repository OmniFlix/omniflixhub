package types

import "C"
import (
	"time"

	"github.com/OmniFlix/omniflixhub/v3/x/itc/exported"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
)

var (
	_ proto.Message      = &Campaign{}
	_ exported.CampaignI = &Campaign{}
)

func NewCampaign(
	id uint64,
	name, description string,
	startTime, endTime time.Time,
	creator, nftDenomId string,
	maxAllowedClaims uint64,
	interaction InteractionType,
	claimType ClaimType,
	tokensPerClaim, totalTokens, availableTokens sdk.Coin,
	nftMintDetails *NFTDetails,
	distribution *Distribution,
) Campaign {
	return Campaign{
		Id:               id,
		Name:             name,
		Description:      description,
		StartTime:        startTime,
		EndTime:          endTime,
		Creator:          creator,
		NftDenomId:       nftDenomId,
		MaxAllowedClaims: maxAllowedClaims,
		Interaction:      interaction,
		ClaimType:        claimType,
		TokensPerClaim:   tokensPerClaim,
		TotalTokens:      totalTokens,
		AvailableTokens:  availableTokens,
		ReceivedNftIds:   []string{},
		NftMintDetails:   nftMintDetails,
		Distribution:     distribution,
		MintCount:        0,
	}
}

func (c Campaign) GetId() uint64 {
	return c.Id
}

func (c Campaign) GetName() string {
	return c.Name
}

func (c Campaign) GetDescription() string {
	return c.Description
}

func (c Campaign) GetStartTime() time.Time {
	return c.StartTime
}

func (c Campaign) GetEndTime() time.Time {
	return c.EndTime
}

func (c Campaign) GetCreator() sdk.AccAddress {
	creator, _ := sdk.AccAddressFromBech32(c.Creator)
	return creator
}

func (c Campaign) GetNftDenomId() string {
	return c.NftDenomId
}

func (c Campaign) GetMaxAllowedClaims() uint64 {
	return c.MaxAllowedClaims
}

func (c Campaign) GetInteraction() string {
	return c.Interaction.String()
}

func (c Campaign) GetClaimType() string {
	return c.ClaimType.String()
}

func (c Campaign) GetTokensPerClaim() sdk.Coin {
	return c.TokensPerClaim
}

func (c Campaign) GetTotalTokens() sdk.Coin {
	return c.TotalTokens
}

func (c Campaign) GetAvailableTokens() sdk.Coin {
	return c.AvailableTokens
}

func (c Campaign) GetNftMintDetails() interface{} {
	return c.NftMintDetails
}

func (c Campaign) GetReceivedNftIds() []string {
	return c.ReceivedNftIds
}

func (c Campaign) GetDistribution() interface{} {
	return c.Distribution
}
