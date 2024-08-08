# oNFT - OmniFlix Non-Fungible Token

The `oNFT` module is a part of the OmniFlix Network and enables the classification and tokenization of assets.

- Assets can be classified under `denoms` (aka `collections` across various ecosystems)
- Tokenize media assets by minting NFTs

`Note:` This module used the [irismod/nft](https://github.com/irismod/nft) repository for initial development and has been modified to meet the requirements of the OmniFlix Network.
### Denom / Collection
A denom is a collection of NFTs. It is a unique identifier for a collection of NFTs.

```protobuf
message Denom {
  option (gogoproto.equal) = true;

  string id                                        = 1;
  string symbol                                    = 2;
  string name                                      = 3;
  string schema                                    = 4;
  string creator                                   = 5;
  string description                               = 6;
  string preview_uri                               = 7 [
    (gogoproto.moretags)   = "yaml:\"preview_uri\"",
    (gogoproto.customname) = "PreviewURI"
  ];
  string uri                                        = 8;
  string uri_hash                                   = 9;
  string data                                       = 10;
  repeated WeightedAddress royalty_receivers        = 11 [
    (gogoproto.moretags) = "yaml:\"royalty_receivers\""
  ];
}
```
## oNFT
oNFT is a non-fungible token that represents a unique asset

```protobuf
message ONFT {
  option (gogoproto.equal)                = true;

  string                    id            = 1;
  Metadata                  metadata      = 2 [(gogoproto.nullable) = false];
  string                    data          = 3;
  string                    owner         = 4;
  bool                      transferable  = 5;
  bool                      extensible    = 6;
  google.protobuf.Timestamp created_at    = 7 [
    (gogoproto.nullable) = false,
    (gogoproto.stdtime)  = true,
    (gogoproto.moretags) = "yaml:\"created_at\""
  ];
  bool                      nsfw          = 8;
  string                    royalty_share = 9 [
    (gogoproto.nullable)   = false,
    (gogoproto.moretags)   = "yaml:\"royalty_share\"",
    (gogoproto.customtype) = "cosmossdk.io/math.LegacyDec"
  ];
}

message Metadata {
  option (gogoproto.equal) = true;

  string name              = 1 [(gogoproto.moretags) = "yaml:\"name\""];
  string description       = 2 [(gogoproto.moretags) = "yaml:\"description\""];
  string media_uri         = 3 [
    (gogoproto.moretags)   = "yaml:\"media_uri\"",
    (gogoproto.customname) = "MediaURI"
  ];
  string preview_uri       = 4 [
    (gogoproto.moretags)   = "yaml:\"preview_uri\"",
    (gogoproto.customname) = "PreviewURI"
  ];
  string uri_hash          = 5;
}
```
### State
The state of the module is expressed by following fields
1. `Collection`: an object contains denom & list of NFTs
2. `Params`: an object contains the parameters of the module


```protobuf
message GenesisState {
  repeated Collection collections = 1 [(gogoproto.nullable) = false];
  Params params = 2 [(gogoproto.nullable) = false];
}

message Collection {
  option (gogoproto.equal) = true;

  Denom         denom      = 1 [(gogoproto.nullable) = false];
  repeated ONFT onfts      = 2 [(gogoproto.customname) = "ONFTs", (gogoproto.nullable) = false];
}

// module params
message Params {
  cosmos.base.v1beta1.Coin     denom_creation_fee = 1 [
    (gogoproto.moretags) = "yaml:\"denom_creation_fee\"",
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Coin",
    (gogoproto.nullable) = false
  ];
}
```



The module supports the following capabilities for classification and tokenization:

- Creation of collections (denoms)
- Minting of NFTs
- Transferring of NFTs
- Burning of NFTs

Various queries are available to get details about denoms/collections, NFTs, and related metadata including but not limited to ownership. Click here to try them out by interacting with the chain.



### 1) Create Denom (Collection)
To create an oNFT denom, you will need to use the "onftd tx onft create" command with the following args and flags:

args:
symbol: denom symbol

flags:
name : name of denom/collection 
description: description for the denom
uri: class uri (ipfs hash of a json file which contains full collections information)
uri-hash: sha256 hash of the above file
preview-uri: display picture url for denom
schema: json schema for additional properties
royalty-receivers: list of weighted addresses that will  receive royalty fees when an NFT is sold
creation-fee: denom creation-fee to create denom

Example:
```
    onftd tx onft create <symbol> \
     --name=<name> \
     --description=<description> \
     --uri=<class-uri> \
     --uri-hash<class-uri-hash> \
     --preview-uri=<preview-uri> \
     --schema=<schema> \
     --royalty-receivers=<address1,weight>,<address2,wight> \ 
     --creation-fee=<creation-fee> \
     --chain-id=<chain-id> \
     --fees=<fee> \
     --from=<key-name>
```
  
### 2) Mint an oNFT

To create an oNFT, you will need to use the "onftd tx onft mint" command with the following flags:

denom-id: the ID of the collection in which you want to mint the NFT
name: name of the NFT
description: description of the NFT
media-uri: the IPFS URI of the NFT (can be a image ipfs url or json file ipfs url that contains all the info of the nft)
preview-uri: the preview URI of the NFT (preview image of the NFT)
data: any additional properties for the NFT in json string format (optional)
recipient: the recipient of the NFT (optional, default is the minter of the NFT)
non-transferable: flag to mint a non-transferable NFT (optional, default is false)
inextensible: flag to mint an inextensible NFT (optional, default is false)
nsfw: flag to mark the NFT as not safe for work (optional, default is false)
royalty-share: the royalty share for the NFT (optional, default is 0.00)

Example:

```
onftd tx onft mint <denom-id> \
--name="NFT name" \
--description="NFT description" \
--media-uri="https://ipfs.io/ipfs/...." \
--uri-hash="uri-hash" \
--preview-uri="https://ipfs.io/ipfs/...." \
--data="{}" \
--recipient="" \
--chain-id=<chain-id> \
--fees=<fee> \
--from=<key-name> 
```

Optional flags:
```
--non-transferable 
--inextensible 
--nsfw 
```

For a royalty share:

```
--royalty-share="0.05" # 5%
```

### 3) Transfer an oNFT

To transfer an oNFT, you will need to use the "onftd tx onft transfer" command with the following flags:

recipient: the recipient's account address
denom-id: the ID of the collection in which the NFT is located
onft-id: the ID of the NFT to be transferred
chain-id: the ID of the blockchain where the transaction will be made (required)
fees: the transaction fees (required)
from: the name of the key to sign the transaction with (required)

Example:

```
onftd tx onft transfer <recipient> <denom-id> <onft-id> \
--chain-id=<chain-id> \
--fees=<fee> \
--from=<key-name>
```

### 4) Burn an oNFT

To burn an oNFT, you will need to use the "onftd tx onft burn" command with the following flags:

denom-id: the ID of the collection in which the NFT is located
onft-id: the ID of the NFT to be burned
chain-id: the ID of the blockchain where the transaction will be made (required)
fees: the transaction fees (required)
from: the name of the key to sign the transaction with (required)

Example:

```
onftd tx onft burn <denom-id> <onft-id> \
--chain-id=<chain-id> \
--fees=<fee> \
--from=<key-name>
```

### Queries
List of queries available for the module:

```protobuf
service Query {
  rpc Collection(QueryCollectionRequest) returns (QueryCollectionResponse) {
    option (google.api.http).get = "/omniflix/onft/v1beta1/collections/{denom_id}";
  }

  rpc IBCCollection(QueryIBCCollectionRequest) returns (QueryCollectionResponse) {
    option (google.api.http).get = "/omniflix/onft/v1beta1/collections/ibc/{hash}";
  }

  rpc Denom(QueryDenomRequest) returns (QueryDenomResponse) {
    option (google.api.http).get = "/omniflix/onft/v1beta1/denoms/{denom_id}";
  }

  rpc IBCDenom(QueryIBCDenomRequest) returns (QueryDenomResponse) {
    option (google.api.http).get = "/omniflix/onft/v1beta1/denoms/ibc/{hash}";
  }

  rpc Denoms(QueryDenomsRequest) returns (QueryDenomsResponse) {
    option (google.api.http).get = "/omniflix/onft/v1beta1/denoms";
  }
  rpc ONFT(QueryONFTRequest) returns (QueryONFTResponse) {
    option (google.api.http).get = "/omniflix/onft/v1beta1/denoms/{denom_id}/onfts/{id}";
  }
  rpc IBCDenomONFT(QueryIBCDenomONFTRequest) returns (QueryONFTResponse) {
    option (google.api.http).get = "/omniflix/onft/v1beta1/denoms/ibc/{hash}/onfts/{id}";
  }
  rpc OwnerONFTs(QueryOwnerONFTsRequest) returns (QueryOwnerONFTsResponse) {
    option (google.api.http).get = "/omniflix/onft/v1beta1/onfts/{denom_id}/{owner}";
  }
  rpc OwnerIBCDenomONFTs(QueryOwnerIBCDenomONFTsRequest) returns (QueryOwnerONFTsResponse) {
    option (google.api.http).get = "/omniflix/onft/v1beta1/onfts/ibc/{hash}/{owner}";
  }
  rpc Supply(QuerySupplyRequest) returns (QuerySupplyResponse) {
    option (google.api.http).get = "/omniflix/onft/v1beta1/denoms/{denom_id}/supply";
  }
  rpc IBCDenomSupply(QueryIBCDenomSupplyRequest) returns (QuerySupplyResponse) {
    option (google.api.http).get = "/omniflix/onft/v1beta1/denoms/ibc/{hash}/supply";
  }
  rpc Params(QueryParamsRequest) returns (QueryParamsResponse) {
    option (google.api.http).get = "/omniflix/onft/v1beta1/params";
  }
}
```
### CLI Queries
  - #### Get List of denoms (collections)
    ```bash
    onftd query onft denoms
    ```
  - #### Get Denom details by it's Id
     ```bash
    onftd query onft denom <denom-id>
    ```    
  - #### Get List of NFTs in a collection
    ```bash
    onftd query onft collection <denom-id>
    ```
  - #### Get Total Count of NFTs in a collection
    ```bash
    onftd query onft supply <denom-id>
    ```
  - #### Get NFT details by it's Id
    ```bash
    onftd query onft asset <denom-id> <nft-id>
    ```
  - #### Get All NFTs owned by an address
    ```bash
    onftd query onft owner <account-address>
    ```
