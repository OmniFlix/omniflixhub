# OmniFlix Hub

OmniFlix Hub is the root chain of the OmniFlix Network. Sovereign chains and DAOs build on top of or connect to the OmniFlix Hub to manage their web2 & web3 media operations to mint, manage, distribute & monetize NFTs enabled with community interactions

OmniFlix Hub blockchain is built using Cosmos-SDK and CometBFT

### Hardware Requirements
 - 6+ core CPUs (recommended: AMD x86_64)
 - 32GB RAM
 - 1TB SSD Storage

### Go Requirement
- go 1.22.5 +
```
sudo rm -rf /usr/local/go
wget -q -O - https://git.io/vQhTU | bash -s -- --remove
wget -q -O - https://git.io/vQhTU | bash -s -- --version 1.21.3
 ```

### Installation

**Linux**

```
git clone https://github.com/Omniflix/omniflixhub.git
cd omniflixhub
git checkout v4.1.1
go mod tidy
make install
```

### Setup

```
# Initialize node
MONIKER=omniflix-node
CHAIN_ID=omniflixhub-1
omniflixhubd init $MONIKER --chain-id $CHAIN_ID

# Download Genesis
curl https://raw.githubusercontent.com/OmniFlix/mainnet/main/omniflixhub-1/genesis.json > ~/.omniflixhub/config/genesis.json

# verify sha256 sum 
# 3c01dd89ae10f3dc247648831ef9e8168afd020946a13055d92a7fe2f50050a0
```
**Update Config** [(omniflixhub-1/config)](https://github.com/OmniFlix/docs/blob/main/guides/mainnet/omniflixhub-1/run-full-node.md#2-update-config)
- update minimum-gas-price as `0.001uflix` in `app.toml`
- updates persistent peers and seeds in config.toml 


### Snapshots
 - [c29r3/cosmos-snapshots](https://github.com/c29r3/cosmos-snapshots) repository for omniflixhub snapshots
 - [polkachu snapshot](https://polkachu.com/tendermint_snapshots/omniflix) 
 - [NodeStake snapshot](https://nodestake.top/omniflix)
   
### OmniFlix Modules
- [oNFT](https://github.com/OmniFlix/omniflixhub/tree/main/x/onft)
- [Marketplace](https://github.com/OmniFlix/omniflixhub/tree/main/marketplace)
- [ITC](https://github.com/OmniFlix/omniflixhub/tree/main/itc)
- [StreamPay](https://github.com/OmniFlix/streampay)
 
### Documentation

- [flixnet-1 guides](https://github.com/OmniFlix/docs/tree/main/guides/testnets/flixnet-1)
- [flixnet-2 guides](https://github.com/OmniFlix/docs/tree/main/guides/testnets/flixnet-2)
- [flixnet-3 guides](https://github.com/OmniFlix/docs/tree/main/guides/testnets/flixnet-3)
- [flixnet-4 guides](https://github.com/OmniFlix/docs/tree/main/guides/testnets/flixnet-4)
- [Mainnet Guides](https://github.com/OmniFlix/docs/tree/main/guides/mainnet)

## Upgrades
 - [upgrade_1](https://github.com/OmniFlix/docs/blob/main/guides/mainnet/omniflixhub-1/upgrades/upgrade_1.md) at block 4175400
 - [v0.10.0](https://github.com/OmniFlix/docs/blob/main/guides/mainnet/omniflixhub-1/upgrades/v0.10.0-upgrade.md) at block 6262420
 - [v0.11.0](https://github.com/OmniFlix/docs/blob/main/guides/mainnet/omniflixhub-1/upgrades/v0.11.0-upgrade.md) at block 7339200
 - [v0.12.x](https://github.com/OmniFlix/docs/blob/main/guides/mainnet/omniflixhub-1/upgrades/v0.12.x-upgrade.md) at block 8054200
 - [v2](https://github.com/OmniFlix/docs/blob/main/guides/mainnet/omniflixhub-1/upgrades/v2-upgrade.md) at block 10428200
 - [v2.1]((https://github.com/OmniFlix/docs/blob/main/guides/mainnet/omniflixhub-1/upgrades/v2.1-upgrade.md)) at block 10678600
 - [v3]((https://github.com/OmniFlix/docs/blob/main/guides/mainnet/omniflixhub-1/upgrades/v3-upgrade.md)) at block 10872800
 - [v3.3.0]((https://github.com/OmniFlix/docs/blob/main/guides/mainnet/omniflixhub-1/upgrades/v3.3.0-upgrade.md)) at block 11140000
  - [v4]((https://github.com/OmniFlix/docs/blob/main/guides/mainnet/omniflixhub-1/upgrades/v4-upgrade.md)) at block 11914000

### Testnets

- [flixnet-1](https://github.com/OmniFlix/testnets/tree/main/flixnet-1)
- [flixnet-2](https://github.com/OmniFlix/testnets/tree/main/flixnet-2)
- [flixnet-3](https://github.com/OmniFlix/testnets/tree/main/flixnet-3)
- [flixnet-4](https://github.com/OmniFlix/testnets/tree/main/flixnet-4)

## Mainnet
- [omniflixhub-1](https://github.com/OmniFlix/mainnet/tree/main/omniflixhub-1)

