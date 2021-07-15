package keeper

import (
	"github.com/omniflix/omniflix/x/omniflix/types"
)

var _ types.QueryServer = Keeper{}
