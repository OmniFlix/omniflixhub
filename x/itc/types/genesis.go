package types

import (
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewGenesisState(campaigns []Campaign,
	claims []Claim,
	nextCampaignNumber uint64,
	params Params,
) *GenesisState {
	return &GenesisState{
		Campaigns:          campaigns,
		Claims:             claims,
		Params:             params,
		NextCampaignNumber: nextCampaignNumber,
	}
}

func (m *GenesisState) ValidateGenesis() error {
	for _, c := range m.Campaigns {
		if c.GetCreator().Empty() {
			return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "missing campaign creator address")
		}
		if err := ValidateCampaign(c); err != nil {
			return err
		}
	}
	for _, claim := range m.Claims {
		if err := ValidateClaim(claim); err != nil {
			return err
		}
	}
	if err := m.Params.ValidateBasic(); err != nil {
		return err
	}
	if m.NextCampaignNumber <= 0 {
		return errorsmod.Wrap(ErrNonPositiveNumber, "must be a number and greater than 0.")
	}
	return nil
}
