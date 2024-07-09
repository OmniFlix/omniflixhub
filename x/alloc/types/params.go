package types

import sdkmath "cosmossdk.io/math"

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
			StakingRewards:      sdkmath.LegacyNewDecWithPrec(60, 2), // 60%
			NftIncentives:       sdkmath.LegacyNewDecWithPrec(15, 2), // 15%
			NodeHostsIncentives: sdkmath.LegacyNewDecWithPrec(5, 2),  // 5%
			DeveloperRewards:    sdkmath.LegacyNewDecWithPrec(15, 2), // 15%
			CommunityPool:       sdkmath.LegacyNewDecWithPrec(5, 2),  // 5%
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
