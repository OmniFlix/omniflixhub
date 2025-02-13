set -ex
MONIKER=test
DENOM=uflix
CHAINID=omniflixhub-test
omniflixhubd version --long

# Setup OmniFlixHub
omniflixhubd init $MONIKER --chain-id $CHAINID --home ./.omniflixhub-test
#sed -i 's#tcp://127.0.0.1:26657#tcp://0.0.0.0:26657#g' ./.omniflixhub-test/config/config.toml
#sed -i "s/\"stake\"/\"$DENOM\"/g" ./.omniflixhu-test/config/genesis.json
#sed -i 's/enable = false/enable = true/g' ./.omniflixhub-test/config/app.toml
omniflixhubd keys add validator --keyring-backend test --home ./.omniflixhub-test

omniflixhubd add-genesis-account validator 100000000000000$DENOM --keyring-backend test --home ./.omniflixhub-test
omniflixhubd gentx validator 900000000$DENOM --keyring-backend test --chain-id $CHAINID --home ./.omniflixhub-test
omniflixhubd collect-gentxs --home ./.omniflixhub-test

#omniflixhubd start --home ./.omniflixhub-test
