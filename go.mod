module github.com/OmniFlix/omniflixhub

go 1.16

require (
	github.com/OmniFlix/marketplace v0.1.0-rc3
	github.com/OmniFlix/onft v0.2.0-rc3
	github.com/cosmos/cosmos-sdk v0.44.3
	github.com/cosmos/ibc-go/v2 v2.0.0
	github.com/gorilla/mux v1.8.0
	github.com/spf13/cast v1.4.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/tendermint/spm v0.1.7
	github.com/tendermint/tendermint v0.34.14
	github.com/tendermint/tm-db v0.6.4
)

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
