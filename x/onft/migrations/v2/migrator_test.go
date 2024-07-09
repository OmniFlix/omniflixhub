package v2_test

import (
	"testing"

	"github.com/OmniFlix/omniflixhub/v5/app/apptesting"
	"github.com/stretchr/testify/suite"

	"github.com/stretchr/testify/require"

	"github.com/OmniFlix/omniflixhub/v5/x/onft/exported"
	v2 "github.com/OmniFlix/omniflixhub/v5/x/onft/migrations/v2"
	"github.com/OmniFlix/omniflixhub/v5/x/onft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MigratorTestSuite struct {
	apptesting.KeeperTestHelper
}

func TestMigratorTestSuite(t *testing.T) {
	suite.Run(t, new(MigratorTestSuite))
}

func (suite *MigratorTestSuite) SetupTest() {
	suite.Setup()
}

type mockSubspace struct {
	ps types.Params
}

func newMockSubspace(ps types.Params) mockSubspace {
	return mockSubspace{ps: ps}
}

func (ms mockSubspace) GetParamSet(ctx context.Context, ps exported.ParamSet) {
	*ps.(*types.Params) = ms.ps
}

func (suite *MigratorTestSuite) TestMigrate() {
	storeKey := suite.App.GetKey(types.StoreKey)
	// tKey := sdk.NewTransientStoreKey("transient_test")
	ctx := suite.Ctx
	cdc := suite.App.AppCodec()
	store := ctx.KVStore(storeKey)
	collections := generateCollectionsData(ctx, storeKey, cdc)

	legacySubspace := newMockSubspace(types.DefaultParams())
	require.NoError(suite.T(), v2.Migrate(ctx, storeKey, legacySubspace, cdc, suite.App.ONFTKeeper))

	var res types.Params
	bz := store.Get(v2.ParamsKey)
	require.NoError(suite.T(), cdc.Unmarshal(bz, &res))
	require.Equal(suite.T(), legacySubspace.ps, res)
	check(suite.T(), ctx, suite.App.ONFTKeeper, collections)
}
