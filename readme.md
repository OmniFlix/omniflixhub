# OmniFlix Hub

**OmniFlix Hub** is a blockchain built using Cosmos SDK and Tendermint and initially created
with [Ignite CLI](https://github.com/ignite/cli).

### Hardware Requirements
 - Quad Core or larger AMD or Intel (amd64) CPU 
 - 16GB RAM
 - 500GB SSD Storage

### Go Requirement
- go 1.19 +
```
sudo rm -rf /usr/local/go
wget -q -O - https://git.io/vQhTU | bash -s -- --remove
wget -q -O - https://git.io/vQhTU | bash -s -- --version 1.19.3
 ```

### Installation

**Linux**

```
git clone https://github.com/OmniFlix/omniflixhub/v2.git
cd omniflixhub
git checkout v0.10.0
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
 - Check [c29r3/cosmos-snapshots](https://github.com/c29r3/cosmos-snapshots) repository for omniflixhub snapshots
   
### OmniFlix Modules
- [oNFT](https://github.com/OmniFlix/onft)
- [Marketplace](https://github.com/OmniFlix/marketplace)
 
### Documentation

- [flixnet-1 guides](https://github.com/OmniFlix/docs/tree/main/guides/testnets/flixnet-1)
- [flixnet-2 guides](https://github.com/OmniFlix/docs/tree/main/guides/testnets/flixnet-2)
- [flixnet-3 guides](https://github.com/OmniFlix/docs/tree/main/guides/testnets/flixnet-3)
- [flixnet-4 guides](https://github.com/OmniFlix/docs/tree/main/guides/testnets/flixnet-4)
- [Mainnet Guides](https://github.com/OmniFlix/docs/tree/main/guides/mainnet)

## Upgrades
 - [upgrade_1](https://github.com/OmniFlix/docs/blob/main/guides/mainnet/omniflixhub-1/upgrades/upgrade_1.md) at block 4175400
 - [v0.10.0](https://github.com/OmniFlix/docs/blob/main/guides/mainnet/omniflixhub-1/upgrades/v0.10.0-upgrade.md) at block 6262420

### Testnets

- [flixnet-1](https://github.com/OmniFlix/testnets/tree/main/flixnet-1)
- [flixnet-2](https://github.com/OmniFlix/testnets/tree/main/flixnet-2)
- [flixnet-3](https://github.com/OmniFlix/testnets/tree/main/flixnet-3)
- [flixnet-4](https://github.com/OmniFlix/testnets/tree/main/flixnet-4)

## Mainnet
- [omniflixhub-1](https://github.com/OmniFlix/mainnet/tree/main/omniflixhub-1)

