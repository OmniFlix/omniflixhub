package types

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys
var (
	KeyDistributionProportions    = []byte("DistributionProportions")
	KeyDeveloperRewardsReceiver   = []byte("DeveloperRewardsReceiver")
	KeyNftIncentivesReceiver      = []byte("NftIncentivesReceiver")
	KeyNodeHostIncentivesReceiver = []byte("NodeHostIncentivesReceiver")
)

// ParamKeyTable ParamTable for module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyDistributionProportions, &p.DistributionProportions, validateDistributionProportions),
		paramtypes.NewParamSetPair(
			KeyDeveloperRewardsReceiver,
			&p.WeightedDeveloperRewardsReceivers,
			validateWeightedAddresses,
		),
		paramtypes.NewParamSetPair(
			KeyNftIncentivesReceiver,
			&p.WeightedNftIncentivesReceivers,
			validateWeightedAddresses,
		),
		paramtypes.NewParamSetPair(
			KeyNodeHostIncentivesReceiver,
			&p.WeightedNodeHostsIncentivesReceivers,
			validateWeightedAddresses,
		),
	}
}

func validateDistributionProportions(i interface{}) error {
	v, ok := i.(DistributionProportions)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.StakingRewards.IsNegative() {
		return errors.New("staking rewards distribution ratio should not be negative")
	}

	if v.NftIncentives.IsNegative() {
		return errors.New("NFT incentives distribution ratio should not be negative")
	}

	if v.NodeHostsIncentives.IsNegative() {
		return errors.New("node hosts incentives distribution ratio should not be negative")
	}

	if v.DeveloperRewards.IsNegative() {
		return errors.New("developer rewards distribution ratio should not be negative")
	}

	if v.CommunityPool.IsNegative() {
		return errors.New("community pool distribution ratio should not be negative")
	}

	totalProportions := v.StakingRewards.Add(v.NftIncentives).Add(v.NodeHostsIncentives).Add(v.DeveloperRewards).Add(v.CommunityPool)

	if !totalProportions.Equal(sdk.NewDec(1)) {
		return errors.New("total distributions ratio should be equal to 100%")
	}

	return nil
}

func validateWeightedAddresses(i interface{}) error {
	v, ok := i.([]WeightedAddress)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// fund community pool when rewards address is empty
	if len(v) == 0 {
		return nil
	}

	weightSum := sdk.NewDec(0)
	for i, w := range v {
		// we allow address to be "" to go to community pool
		if w.Address != "" {
			_, err := sdk.AccAddressFromBech32(w.Address)
			if err != nil {
				return fmt.Errorf("invalid address at %dth", i)
			}
		}
		if !w.Weight.IsPositive() {
			return fmt.Errorf("non-positive weight at %dth", i)
		}
		if w.Weight.GT(sdk.NewDec(1)) {
			return fmt.Errorf("more than 1 weight at %dth", i)
		}
		weightSum = weightSum.Add(w.Weight)
	}

	if !weightSum.Equal(sdk.NewDec(1)) {
		return fmt.Errorf("invalid weight sum: %s", weightSum.String())
	}

	return nil
}
