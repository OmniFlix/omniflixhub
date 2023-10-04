package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewParams(
	distrProportions DistributionProportions,
	weightedDevRewardsReceivers []WeightedAddress,
	weightedNftIncentivesReceivers []WeightedAddress,
	weightedNodeHostIncentivesReceivers []WeightedAddress,
) Params {
	return Params{
		DistributionProportions:              distrProportions,
		WeightedDeveloperRewardsReceivers:    weightedDevRewardsReceivers,
		WeightedNftIncentivesReceivers:       weightedNftIncentivesReceivers,
		WeightedNodeHostsIncentivesReceivers: weightedNodeHostIncentivesReceivers,
	}
}

// DefaultParams default module parameters
func DefaultParams() Params {
	return Params{
		DistributionProportions: DistributionProportions{
			StakingRewards:      sdk.NewDecWithPrec(60, 2), // 60%
			NftIncentives:       sdk.NewDecWithPrec(15, 2), // 15%
			NodeHostsIncentives: sdk.NewDecWithPrec(5, 2),  // 5%
			DeveloperRewards:    sdk.NewDecWithPrec(15, 2), // 15%
			CommunityPool:       sdk.NewDecWithPrec(5, 2),  // 5%
		},
		WeightedDeveloperRewardsReceivers:    []WeightedAddress(nil),
		WeightedNftIncentivesReceivers:       []WeightedAddress(nil),
		WeightedNodeHostsIncentivesReceivers: []WeightedAddress(nil),
	}
}

// Validate validate params
func (p Params) Validate() error {
	if err := validateDistributionProportions(p.DistributionProportions); err != nil {
		return err
	}
	err := validateWeightedAddresses(p.WeightedDeveloperRewardsReceivers)
	return err
}
