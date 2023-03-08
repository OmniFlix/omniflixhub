#!/bin/bash

set -eo pipefail

# get protoc executions
#go get github.com/regen-network/cosmos-proto/protoc-gen-gocosmos 2>/dev/null
# get cosmos sdk from github
#go get github.com/cosmos/cosmos-sdk 2>/dev/null

# Get the path of the cosmos-sdk repo from go/pkg/mod
 # cosmos_sdk_dir=$(go list -f '{{ .Dir }}' -m github.com/cosmos/cosmos-sdk)
# echo $cosmos_sdk_dir;
proto_dirs=$(find ./proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  protoc \
  -I "proto" \
  -I "third_party/proto" \
  --gocosmos_out=plugins=interfacetype+grpc,\
Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types:. \
  $(find "${dir}" -maxdepth 1 -name '*.proto')

  # command to generate gRPC gateway (*.pb.gw.go in respective modules) files
  protoc \
  -I "proto" \
  -I "third_party/proto" \
  --grpc-gateway_out=logtostderr=true:. \
  $(find "${dir}" -maxdepth 1 -name '*.proto')
done

# move proto files to the right places
cp -r github.com/OmniFlix/omniflixhub/* ./
rm -rf github.com
