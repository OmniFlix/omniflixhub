#!/usr/bin/env bash

set -eo pipefail

mkdir -p ./tmp-swagger-gen

# Get the paths used repos from go/pkg/mod
cosmos_sdk_dir=$(go list -f '{{ .Dir }}' -m github.com/cosmos/cosmos-sdk)
ibc_go=$(go list -f '{{ .Dir }}' -m github.com/cosmos/ibc-go/v8)
streampay=$(go list -f '{{ .Dir }}' -m github.com/OmniFlix/streampay/v2)
echo "$cosmos_sdk_dir"/proto
cd proto
proto_dirs=$(find ./OmniFlix ./osmosis "$ibc_go"/proto "$streampay"/proto "$cosmos_sdk_dir"/proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  # generate swagger files (filter query files)
  query_file=$(find "${dir}" -maxdepth 1 \( -name 'query.proto' -o -name 'service.proto' \))
  if [[ ! -z "$query_file" ]]; then
    buf generate --template buf.gen.swagger.yaml $query_file
  fi
done

# come back to main folder
cd ..

# Fix circular definition in cosmos/tx/v1beta1/service.swagger.json
jq 'del(.definitions["cosmos.tx.v1beta1.ModeInfo.Multi"].properties.mode_infos.items["$ref"])' ./tmp-swagger-gen/cosmos/tx/v1beta1/service.swagger.json > ./tmp-swagger-gen/cosmos/tx/v1beta1/fixed-service.swagger.json
jq 'del(.definitions["cosmos.autocli."].properties.mode_infos.items["$ref"])' ./tmp-swagger-gen/cosmos/tx/v1beta1/fixed-service.swagger.json > ./tmp-swagger-gen/cosmos/tx/v1beta1/fixed-service2.swagger.json

rm ./tmp-swagger-gen/cosmos/tx/v1beta1/service.swagger.json

# Tag everything as "gRPC Gateway API"
perl -i -pe 's/"(Query|Service)"/"gRPC Gateway API"/' $(find ./tmp-swagger-gen -name '*.swagger.json' -print0 | xargs -0)

# Convert all *.swagger.json files into a single folder _all
files=$(find ./tmp-swagger-gen -name '*.swagger.json' -print0 | xargs -0)
mkdir -p ./tmp-swagger-gen/_all
counter=0
for f in $files; do
  echo "[+] $f"

  if [[ "$f" =~ "cosmos" ]]; then
    cp $f ./tmp-swagger-gen/_all/cosmos-$counter.json
  elif [[ "$f" =~ "omniflix" ]]; then
    cp $f ./tmp-swagger-gen/_all/omniflix-$counter.json
  elif [[ "$f" =~ "streampay" ]]; then
    cp $f ./tmp-swagger-gen/_all/streampay-$counter.json
  else
    cp $f ./tmp-swagger-gen/_all/other-$counter.json
  fi
  ((counter++))
done

# merges all the above into FINAL.json
python3 ./scripts/merge_protoc.py

# Makes a swagger temp file with reference pointers
swagger-combine ./tmp-swagger-gen/_all/FINAL.json -o ./docs/_tmp_swagger.yaml -f yaml --continueOnConflictingPaths --includeDefinitions

# extends out the *ref instances to their full value
swagger-merger --input ./docs/_tmp_swagger.yaml -o ./docs/swagger.yaml

# Derive openapi from swagger docs
swagger2openapi --patch ./docs/swagger.yaml --outfile ./docs/static/openapi.yml --yaml

# clean swagger tmp files
rm ./docs/_tmp_swagger.yaml
rm -rf ./tmp-swagger-gen