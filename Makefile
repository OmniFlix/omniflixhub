PACKAGES=$(shell go list ./... | grep -v '/simulation')

VERSION := $(shell echo $(shell git describe --tags --always) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
COSMOS_SDK := $(shell grep -i cosmos-sdk go.mod | awk '{print $$2}')

build_tags := $(strip netgo $(build_tags))

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=OmniFlix \
	-X github.com/cosmos/cosmos-sdk/version.AppName=omniflixd \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
	-X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags),cosmos-sdk $(COSMOS_SDK)"

BUILD_FLAGS := -ldflags '$(ldflags)'

all: go.sum install

create-wallet:
	omniflixd keys add mywallet

init-chain: create-wallet
	rm -rf ~/.omniflix
	omniflixd init omniflix-node  --chain-id "omniflix-test" --stake-denom uflix
	omniflixd add-genesis-account $(shell omniflixd keys show mywallet -a) 1000000000000uflix
	omniflixd gentx --name=mywallet --amount 100000000uflix
	omniflixd collect-gentxs

install: go.sum
		go build $(BUILD_FLAGS) -o ${GOPATH}/bin/omniflixd ./cmd/omniflixd
build:
		go build $(BUILD_FLAGS) -o ${GOPATH}/bin/omniflixd ./cmd/omniflixd

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify

lint:
	@echo "--> Running linter"
	@golangci-lint run
	@go mod verify
