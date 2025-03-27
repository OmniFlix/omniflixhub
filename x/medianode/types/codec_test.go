package types

import (
	"testing"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/suite"
)

type CodecTestSuite struct {
	suite.Suite
}

func TestCodecSuite(t *testing.T) {
	suite.Run(t, new(CodecTestSuite))
}

func (suite *CodecTestSuite) TestRegisterInterfaces() {
	registry := codectypes.NewInterfaceRegistry()
	registry.RegisterInterface(sdk.MsgInterfaceProtoName, (*sdk.Msg)(nil))
	RegisterInterfaces(registry)

	impls := registry.ListImplementations(sdk.MsgInterfaceProtoName)
	suite.Require().Equal(7, len(impls))
	suite.Require().ElementsMatch([]string{
		"/OmniFlix.medianode.v1.MsgRegisterMediaNode",
		"/OmniFlix.medianode.v1.MsgDepositMediaNode",
		"/OmniFlix.medianode.v1.MsgUpdateMediaNode",
		"/OmniFlix.medianode.v1.MsgLeaseMediaNode",
		"/OmniFlix.medianode.v1.MsgExtendLease",
		"/OmniFlix.medianode.v1.MsgCancelLease",
		"/OmniFlix.medianode.v1.MsgCloseMediaNode",
	}, impls)
}
