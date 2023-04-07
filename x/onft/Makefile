APP_NAME = onft
DAEMON_NAME = onftd
LEDGER_ENABLED ?= true

PACKAGES=$(shell go list ./... | grep -v '/simulation')
VERSION := $(shell echo $(shell git describe --tags --always) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
COSMOS_SDK := $(shell grep -i cosmos-sdk go.mod | awk '{print $$2}')

build_tags = netgo,
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags+=ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags+=ledger
      endif
    endif
  endif
endif
build_tags := $(strip $(build_tags))

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=${APP_NAME} \
	-X github.com/cosmos/cosmos-sdk/version.AppName=${DAEMON_NAME} \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
	-X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags),cosmos-sdk $(COSMOS_SDK)"

BUILD_FLAGS := -ldflags '$(ldflags)'

all: go.sum install

install: go.sum
		go install $(BUILD_FLAGS) ./cmd/onftd/
build:
		go build $(BUILD_FLAGS) -o ${GOPATH}/bin/${DAEMON_NAME} ./cmd/onftd/

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify

lint:
	@echo "--> Running linter"
	@golangci-lint run
	@go mod verify

start-test-chain:
	rm  ~/.onft/config/genesis.json
	rm -rf  ~/.onft/config/gentx
	onftd unsafe-reset-all
	onftd init onft-node  --chain-id "onft-test-1"
	onftd keys add validator
	onftd add-genesis-account `onftd keys show validator -a` 100000000uflix
	onftd gentx validator 1000000uflix --moniker "validator-1" --chain-id "onft-test-1"
	sed -i "s/\"stake\"/\"uflix\"/g" ~/.onft//config/genesis.json
	onftd collect-gentxs
	onftd validate-genesis
	onftd start

########################################
### Testing
SIMAPP = ./app

PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
PACKAGES_UNITTEST=$(shell go list ./... | grep -v '/simulation' | grep -v '/cli_test')

test: test-unit

test-unit:
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' -ldflags '$(ldflags)' ${PACKAGES_UNITTEST}

test-sim-nondeterminism:
	@echo "Running non-determinism test..."
	@go test -mod=readonly $(SIMAPP) -run TestAppStateDeterminism -Enabled=true \
		-NumBlocks=100 -BlockSize=200 -Commit=true -Period=0 -v -timeout 24h

test-sim-nondeterminism-fast:
	@echo "Running non-determinism test..."
	@go test -mod=readonly $(SIMAPP) -run TestAppStateDeterminism -Enabled=true \
		-NumBlocks=10 -BlockSize=200 -Commit=true -Period=0 -v -timeout 24h

test-sim-custom-genesis-fast:
	@echo "Running custom genesis simulation..."
	@echo "By default, $(shell pwd)/testdata/genesis.json will be used."
	@go test -mod=readonly $(SIMAPP) -run TestFullAppSimulation -Genesis=$(shell pwd)/testdata/genesis.json \
		-Enabled=true -NumBlocks=10 -BlockSize=200 -Commit=true -Seed=99 -Period=5 -v -timeout 24h