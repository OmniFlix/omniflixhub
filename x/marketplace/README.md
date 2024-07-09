# x/marketplace

The `marketplace` module allows users to list NFTs with a fixed price or through timed auctions. Buyers can purchase NFTs from the marketplace by paying the listed price or by participating in an auction. The `marketplace` module supports different types of listings:

- Fixed Price Listing
- Timed Auction
### Fixed Price Listing

In a fixed price listing, buyers must pay the listed amount to acquire the NFT. NFT owners can list their NFT with any allowed token and specify split shares between different addresses and the percentage of revenue each address receives.

```go
message Listing {
  string                   id       = 1;
  string                   nft_id   = 2 [(gogoproto.moretags) = "yaml:\"nft_id\""];
  string                   denom_id = 3 [(gogoproto.moretags) = "yaml:\"denom_id\""];
  cosmos.base.v1beta1.Coin price    = 4 [
    (gogoproto.nullable)     = false,
    (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coin"
  ];
  string                   owner    = 5 [(gogoproto.moretags) = "yaml:\"owner\""];
  repeated WeightedAddress split_shares  = 6 [
    (gogoproto.nullable) = false,
    (gogoproto.moretags)   = "yaml:\"split_shares\""
  ];
}

message WeightedAddress {
  option (gogoproto.equal) = true;

  string address           = 1 [(gogoproto.moretags) = "yaml:\"address\""];
  string weight            = 2 [
    (gogoproto.moretags)   = "yaml:\"weight\"",
    (gogoproto.customtype) = "cosmossdk.io/math.LegacyDec",
    (gogoproto.nullable)   = false
  ];
}
```

## Timed Auction

- Timed auction require buyers to bid with a minimum amount.
- Only one bid is allowed at a time in these auctions.
- When a new bid is placed, the previous bid amount will be returned to the bidder.
- Auction will end at the end time and the highest bidder will be the winner.
- if no bids were placed on the auction, it will be closed at end time.`


```go
message AuctionListing {
  uint64                    id                   = 1;
  string                    nft_id               = 2 [(gogoproto.moretags) = "yaml:\"nft_id\""];
  string                    denom_id             = 3 [(gogoproto.moretags) = "yaml:\"denom_id\""];
  cosmos.base.v1beta1.Coin  start_price          = 4 [
    (gogoproto.nullable)     = false,
    (gogoproto.moretags)     = "yaml:\"start_price\"",
    (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coin"
  ];
  google.protobuf.Timestamp start_time           = 5 [
    (gogoproto.nullable) = false,
    (gogoproto.stdtime)  = true,
    (gogoproto.moretags) = "yaml:\"start_time\""
  ];
  google.protobuf.Timestamp end_time             = 6 [
    (gogoproto.stdtime)  = true,
    (gogoproto.moretags) = "yaml:\"end_time\""
  ];
  string                    owner                = 7;
  string                    increment_percentage = 8 [
    (gogoproto.moretags)   = "yaml:\"increment_percentage\"",
    (gogoproto.customtype) = "cosmossdk.io/math.LegacyDec",
    (gogoproto.nullable)   = false
  ];
  repeated string           whitelist_accounts   = 9 [
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"whitelist_accounts\""
  ];
  repeated WeightedAddress  split_shares         = 10 [
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"split_shares\""
  ];
}
```
## Fees and Distribution

Whenever an NFT is bought or an auction is concluded, a certain percentage of the sale amount is collected as a commission. This commission is then distributed among different parties based on the distribution parameters that have been set.
- **Staking distribution**: The percentage of the fee that is distributed to stakers.
- **Community pool distribution**: The percentage of the fee that is added to the community pool.

## State
The state of the module is expressed by following fields

1. `listings`: A list of NFT listings that are available in the marketplace.
2. `ListingCount`: The total count of NFT listings.
3. `params`: Marketplace module parameters.
4. `auctions`: A list of active auctions in the marketplace.
5. `bids`: A list of bids made on auctions.
6. `next_auction_number`: The number to be assigned to the next auction that is created.

```go
message GenesisState {
  // NFTs that are listed in marketplace
  repeated Listing        listings            = 1 [(gogoproto.nullable) = false];
  uint64                  ListingCount        = 2;
  Params                  params              = 3 [(gogoproto.nullable) = false];
  repeated AuctionListing auctions            = 4 [(gogoproto.nullable) = false];
  repeated Bid            bids                = 5 [(gogoproto.nullable) = false];
  uint64                  next_auction_number = 6;
}
```
### Module parameters
The Marketplace module includes the following parameters:

- **Sale commission**: A decimal that represents the percentage of the sale price that is collected as a fee. This fee is distributed to various parties according to the configured distribution parameters.
- **Distribution**: A data structure that represents the distribution of the sale commission. It includes two fields: the staking distribution and the community pool distribution. These fields are represented as decimals and represent the percentage of the sale commission that is distributed to each party.
- **Bid close duration**: A duration that represents the amount of time after the last bid was placed before the auction is closed. If no bids were placed on the auction, it will be closed after this duration.

## Transactions / Messages
### List NFT
`MsgListNFT` can be submitted by any account to list nft on marketplace.
```go
message MsgListNFT {
  string id = 1;
  string nft_id = 2;
  string denom_id = 3;
  cosmos.base.v1beta1.Coin price = 4 [
    (gogoproto.nullable) = false,
    (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coin"
  ];
  string owner = 5;
  repeated WeightedAddress split_shares = 6 [
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"split_shares\""
  ];
}
```

### Edit Listing
`MsgEditListing` can be submitted by owner of the nft to update price of the nft in `Listing`
```go
message MsgEditListing {
  string id = 1;
  cosmos.base.v1beta1.Coin price = 2 [
    (gogoproto.nullable) = false,
    (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coin"
  ];
  string owner = 3;
}
```

### Remove Listing
`MsgDeListNFT` can be used by owner of the nft to remove the `Listing` from marketplace.
```
message MsgDeListNFT {
  string id = 1;
  string owner = 2;
}
```

### Buy Listed NFT
`MsgBuyNFT` can be used by anyone to buy nft from marketplace.
```go
message MsgBuyNFT {
  string id = 1;
  cosmos.base.v1beta1.Coin price = 2 [
    (gogoproto.nullable) = false,
    (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coin"
  ];
  string buyer = 3;
}
```
### Create Timed Auction
`MsgCreateAuction` can be submitted by any account to create a `Timed Action`
```go
message MsgCreateAuction {
  string nft_id = 1;
  string denom_id = 2;
  google.protobuf.Timestamp start_time = 3 [
    (gogoproto.nullable) = false,
    (gogoproto.stdtime) = true,
    (gogoproto.moretags) = "yaml:\"start_time\""
  ];
  cosmos.base.v1beta1.Coin start_price = 4 [
    (gogoproto.nullable) = false,
    (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coin",
    (gogoproto.moretags) = "yaml:\"start_price\""
  ];
  google.protobuf.Duration duration = 5 [(gogoproto.stdduration) = true];
  string increment_percentage = 6 [
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"increment_percentage\"",
    (gogoproto.customtype) = "cosmossdk.io/math.LegacyDec"
  ];
  repeated string whitelist_accounts = 7 [
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"whitelist_accounts\""
  ];
  repeated WeightedAddress split_shares = 8 [
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"split_shares\""
  ];
  string owner = 9;
}
```

### Cancel Timed Auction
`MsgCancelAuction` can be submitted by auction creator to cancel `Auction` before start time.
```go
message MsgCancelAuction {
  uint64 auction_id = 1 [(gogoproto.moretags) = "yaml:\"auction_id\""];
  string owner = 2;
}
```
`MsgPlaceBid` can be submitted by any account to bid on a `Auction`.
### Place Bid on Auction
```go
message MsgPlaceBid {
  uint64 auction_id = 1 [(gogoproto.moretags) = "yaml:\"auction_id\""];
  cosmos.base.v1beta1.Coin amount = 2 [
    (gogoproto.nullable) = false,
    (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coin"
  ];
  string bidder = 3;
}
```

## Queries

```go
service Query {
  // Params queries params of the marketplace module.
  rpc Params(QueryParamsRequest) returns (QueryParamsResponse) {
    option (google.api.http).get = "/omniflix/marketplace/v1beta1/params";
  }

  rpc Listings(QueryListingsRequest) returns (QueryListingsResponse) {
    option (google.api.http).get = "/omniflix/marketplace/v1beta1/listings";
  }

  rpc Listing(QueryListingRequest) returns (QueryListingResponse) {
    option (google.api.http).get = "/omniflix/marketplace/v1beta1/listings/{id}";
  }

  rpc ListingsByOwner(QueryListingsByOwnerRequest) returns (QueryListingsByOwnerResponse) {
    option (google.api.http).get = "/omniflix/marketplace/v1beta1/listings-by-owner/{owner}";
  }

  rpc ListingsByPriceDenom(QueryListingsByPriceDenomRequest) returns (QueryListingsByPriceDenomResponse) {
    option (google.api.http).get = "/omniflix/marketplace/v1beta1/listings-by-price-denom/{price_denom}";
  }

  rpc ListingByNftId(QueryListingByNFTIDRequest) returns (QueryListingResponse) {
    option (google.api.http).get = "/omniflix/marketplace/v1beta1/listing-by-nft/{nft_id}";
  }

```
- Query Listings
  ```shell
   omniflixhubd q marketplace listings [Flags]
  ```
- Query Listing Details
  ```shell
   omniflixhubd q marketplace listing [listingId] [Flags]
  ```
- Query listings by owner
  ```shell
   omniflixhubd q marketplace listings-by-owner  [owner] [Flags]
  ```
- Query Auctions
  ```shell
   omniflixhubd q marketplace auctions [Flags]
  ```
- Query Auction
  ```shell
   omniflixhubd q marketplace auction <auction-id> [Flags]
  ```
- Query Auction Bid
   ```shell
    omniflixhubd q marketplace bid <auction-id> [Flags]
   ```
