set -ex
MONIKER=test
DENOM=uflix
CHAINID=omniflixhub
ACCKEY=omniflix17yj7k0sey3z3690e4g0hg9q7pcq085u2u6xp2y
omniflixhubd version --long

# Setup OmniFlixHub
omniflixhubd init $MONIKER --chain-id $CHAINID
sed -i 's#tcp://127.0.0.1:26657#tcp://0.0.0.0:26657#g' ~/.omniflixhub/config/config.toml
sed -i "s/\"stake\"/\"$DENOM\"/g" ~/.omniflixhub/config/genesis.json
sed -i 's/enable = false/enable = true/g' ~/.omniflixhub/config/app.toml
omniflixhubd keys add validator --keyring-backend test

omniflixhubd add-genesis-account $(omniflixhubd keys --keyring-backend test show validator -a) 100000000000$DENOM
omniflixhubd add-genesis-account $ACCKEY 100000000000$DENOM
omniflixhubd gentx validator 900000000$DENOM --keyring-backend test --chain-id $CHAINID
omniflixhubd collect-gentxs

omniflixhubd start