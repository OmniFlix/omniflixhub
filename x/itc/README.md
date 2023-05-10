# ITC (Interact To Claim)
## Abstract

The ITC module is designed to enable any account to create campaigns with different types of claims and interactions. Campaigns can be used to distribute fungible tokens or NFTs to users who interact with them using NFTs.

After successfully interacting with a campaign, the claimer will receive the following rewards based on the claim type:

- If the claim type is `CLAIM_TYPE_FT`, the per-claim amount will be distributed to the claimer according to the distribution type.
- If the claim type is `CLAIM_TYPE_NFT`, an NFT will be minted and sent to the claimer address.
- If the claim type is `CLAIM_TYPE_FT_AND_NFT`, both an amount and an NFT will be sent to the claimer address.

### Campaign Types

There are three campaign types that can be created using the ITC module:

- Fungible token claim campaign (`CLAIM_TYPE_FT`): Users can claim fungible tokens by interacting with the campaign using NFTs.
- Non-fungible token claim campaign (`CLAIM_TYPE_NFT`): Users can claim NFTs by interacting with the campaign using NFTs.
- Both fungible and non-fungible token claim campaign (`CLAIM_TYPE_FT_AND_NFT`): Users can claim both fungible tokens and NFTs by interacting with the campaign using NFTs.

### Interaction Types

There are three interaction types that can be used in the ITC module:

- Burn NFT (`INTERACTION_TYPE_BURN`): The NFT will be burned after the claim is made.
- Transfer NFT (`INTERACTION_TYPE_TRANSFER`): The NFT will be transferred to the campaign owner after the claim is made.
- Hold NFT (`INTERACTION_TYPE_HOLD`): The user must hold ownership of the NFT to be able to make a claim.


## State
### Module State
The initial state of the module. It contains information about all campaigns, claims, and parameters of the module.
```go
message GenesisState {
  repeated Campaign    campaigns = 1 [(gogoproto.nullable) = false];
  uint64               next_campaign_number = 2;
  repeated Claim       claims = 3 [(gogoproto.nullable) = false];
  Params               params = 4 [(gogoproto.nullable) = false];
}
```
### Campaign
A campaign created by a user to distribute tokens to other users who perform certain interactions. It has a name, description, start and end time, maximum number of allowed claims, interaction type, claim type, tokens per claim, total tokens, available tokens, received NFT IDs, NFT mint details, distribution, and creator.
```go
message Campaign {
    uint64 id = 1;
    string name = 2;
    string description = 3;
    google.protobuf.Timestamp start_time = 4 [
        (gogoproto.nullable) = false,
        (gogoproto.stdtime) = true,
        (gogoproto.moretags) = "yaml:\"start_time\""
    ];
    google.protobuf.Timestamp end_time = 5 [
        (gogoproto.nullable) = false,
        (gogoproto.stdtime) = true,
        (gogoproto.moretags) = "yaml:\"end_time\""
    ];
    string creator = 6;
    string nft_denom_id = 7 [(gogoproto.moretags) = "yaml:\"nft_denom_id\""];
    uint64 max_allowed_claims = 8
        [(gogoproto.moretags) = "yaml:\"max_allowed_claims\""];
    InteractionType interaction = 9;
    ClaimType claim_type = 10;
    cosmos.base.v1beta1.Coin tokens_per_claim = 11 [
        (gogoproto.nullable) = false,
        (gogoproto.moretags) = "yaml:\"tokens_per_claim\""
    ];
    cosmos.base.v1beta1.Coin total_tokens = 12 [
        (gogoproto.nullable) = false,
        (gogoproto.moretags) = "yaml:\"total_tokens\""
    ];
    cosmos.base.v1beta1.Coin available_tokens = 13 [
        (gogoproto.nullable) = false,
        (gogoproto.moretags) = "yaml:\"available_tokens\""
    ];
    repeated string received_nft_ids = 14
        [(gogoproto.moretags) = "yaml:\"received_nft_ids\""];
    NFTDetails nft_mint_details = 15
        [(gogoproto.moretags) = "yaml:\"nft_mint_details\""];
    Distribution distribution = 16
        [(gogoproto.moretags) = "yaml:\"distribution\""];
}

```
### Claim
A claim made by a user for a campaign by providing their address, NFT ID, and interaction type.
```go
message Claim {
  uint64     campaign_id = 1;
  string     address = 2;
  string     nft_id = 3;
  InteractionType interaction = 4;
}
```
### Parameters
The parameters of the module, which include the maximum campaign duration and creation fee.
```go
message Params {
  google.protobuf.Duration max_campaign_duration = 1 [
    (gogoproto.nullable) = false,
    (gogoproto.stdduration) = true,
    (gogoproto.moretags) = "yaml:\"max_campaign_duration\""
  ];
  cosmos.base.v1beta1.Coin creation_fee = 2 [
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"creation_fee\""
  ];
}
```

## Messages
### CreateCampaign
`MsgCreateCampaign` can be used by any account to create a new `Campaign`
```go
message MsgCreateCampaign {
    string name = 1;
    string description = 2;
    InteractionType interaction = 3;
    ClaimType claim_type = 4 [(gogoproto.moretags) = "yaml:\"claim_type\""];
    string nft_denom_id = 5 [(gogoproto.moretags) = "yaml:\"nft_denom_id\""];
    cosmos.base.v1beta1.Coin tokens_per_claim = 6 [
        (gogoproto.nullable) = false,
        (gogoproto.moretags) = "yaml:\"tokens_per_claim\""
    ];
    uint64 max_allowed_claims = 7
        [(gogoproto.moretags) = "yaml:\"max_allowed_claims\""];
    cosmos.base.v1beta1.Coin deposit = 8 [
        (gogoproto.nullable) = false,
        (gogoproto.moretags) = "yaml:\"deposit\""
    ];
    NFTDetails nft_mint_details = 9
        [(gogoproto.moretags) = "yaml:\"nft_details\""];
    google.protobuf.Timestamp start_time = 10 [
        (gogoproto.nullable) = false,
        (gogoproto.stdtime) = true,
        (gogoproto.moretags) = "yaml:\"start_time\""
    ];

    google.protobuf.Duration duration = 11
        [(gogoproto.nullable) = false, (gogoproto.stdduration) = true];

    Distribution distribution = 12;
    string creator = 13;
    cosmos.base.v1beta1.Coin creation_fee = 14 [
        (gogoproto.nullable) = false,
        (gogoproto.moretags) = "yaml:\"creation_fee\""
    ];
}
```
### CancelCampaign
`MsgCancelCampaign` can be used by the creator to cancel the upcoming `Campaign`.
```go
message MsgCancelCampaign {
    uint64 campaign_id = 1;
    string creator = 2;
}
```

### DepositCampaign
`MsgDepositCampaign` can be used by creator to add funds to the `Campaign`.
```go
message MsgDepositCampaign {
    uint64 campaign_id = 1;
    cosmos.base.v1beta1.Coin amount = 2 [(gogoproto.nullable) = false];
    string depositor = 3;
}
```
### Claim
`MsgClaim` can be used by any account to claim amount/nft from the `Campaign` by interacting with nft
```go
message MsgClaim {
    uint64 campaign_id = 1;
    string nft_id = 2;
    InteractionType interaction = 3;
    string claimer = 4;
}
```

## Transactions

### Create a new campaign

This transaction creates a new campaign with the specified parameters:

```shell
omniflixhubd tx itc create-campaign \
  --name="campaign name" \
  --description="campaign description" \
  --start-time="2023-05-10T10:20:00Z" \
  --duration="3600s" \
  --claim-type=fungible \
  --max-allowed-claims=10 \
  --tokens-per-claim=10000000uflix \
  --deposit=100000000uflix  \
  --interaction-type=transfer \
  --nft-denom-id=nftdenomid \
  --distribution-type stream \
  --stream-duration "600s" \
  --creation-fee 10000000uflix \
  --from=wallet
```

### Cancel campaign

This transaction cancels a campaign with the specified campaign ID:

```shell
omniflixhubd tx itc cancel-campaign <campaign-id> --from wallet
```

### Deposit amount in campaign

This transaction deposits a specified amount into a campaign with the specified campaign ID:

```shell
omniflixhubd tx itc campaign-deposit <campaign-id> --amount 1000000uflix --from wallet
```

### Claim from campaign

This transaction claims a token or NFT from the specified campaign:

```shell
omniflixhubd tx itc claim <campaign-id> \
  --nft-id="nft-id" \
  --interaction-type=<interaction-type> \
  --from wallet
```

**Note:** Replace the values enclosed in `<` and `>` with the actual values.


## Queries
The ITC module provides several queries to fetch information related to campaigns, claims, and module parameters.
```go
service Query {
  rpc Params(QueryParamsRequest) returns (QueryParamsResponse) {
    option (google.api.http).get = "/omniflix/itc/v1/params";
  }

  rpc Campaigns(QueryCampaignsRequest) returns (QueryCampaignsResponse) {
    option (google.api.http).get = "/omniflix/itc/v1/campaigns";
  }

  rpc Campaign(QueryCampaignRequest) returns (QueryCampaignResponse) {
    option (google.api.http).get = "/omniflix/itc/v1/campaigns/{campaign_id}";
  }

  rpc Claims(QueryClaimsRequest) returns (QueryClaimsResponse) {
    option (google.api.http).get = "/omniflix/itc/v1/claims";
  }
}
```

The available queries are:
- `Params`: Returns the current module parameters.
- `Campaigns`: Returns a list of all campaigns.
- `Campaign`: Returns information about a specific campaign by its ID.
- `Claims`: Returns a list of all claims for a specific campaign.

To execute the queries, you can use the following commands:

### Query all campaigns

```shell
omniflixhubd q itc campaigns
```

### Query campaign by ID

```shell
omniflixhubd q itc campaign <campaign-id>
```

### Query claims of a campaign

```shell
omniflixhubd q itc claims <campaign-id>
```

### Query module parameters

```shell
omniflixhubd q itc params
```