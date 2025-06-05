# x/medianode

The `medianode` module enables users to register, manage, and lease media nodes on the OmniFlix Hub. Media nodes are computing resources that can be leased by users for various media-related tasks.

## Abstract

The medianode module allows users to:
- Register new media nodes with specific hardware specifications and pricing
- Update media node information and specifications
- Lease media nodes for a specified duration
- Extend existing leases
- Cancel leases
- Deposit additional funds to media nodes
- Close media nodes and withdraw deposits

## State

### Module State
The initial state of the module contains information about all media nodes, leases, and parameters.
```protobuf
message GenesisState {
  // params defines all the parameters of the module
  Params params                    = 1 [(gogoproto.nullable) = false];
  
  // media_nodes is the list of registered media nodes
  repeated MediaNode nodes   = 2 [(gogoproto.nullable) = false];
  
  // leases is the list of active leases
  repeated Lease leases            = 3 [(gogoproto.nullable) = false];
  
  // next_medianode_id is the current usable media node ID
  uint64 node_counter         = 4;
}
```

### MediaNode
A media node represents a computing resource that can be leased. It contains information about the node's specifications, status, and pricing.
```protobuf
message MediaNode {
  string id                                       = 1;
  string url                                      = 2;
  HardwareSpecs hardware_specs                    = 3 [(gogoproto.nullable) = false, (gogoproto.moretags) = "yaml:\"hardware_specs\""];
  string owner                                    = 4;
  cosmos.base.v1beta1.Coin price_per_hour          = 5 [(gogoproto.nullable) = false, (gogoproto.moretags) = "yaml:\"price_per_hour\""];
  Status status                                   = 6; 
  bool leased                                     = 7;
  google.protobuf.Timestamp registered_at         = 8 [(gogoproto.stdtime) = true, (gogoproto.nullable) = false];
  repeated Deposit deposits                       = 9;
  google.protobuf.Timestamp closed_at             = 10 [(gogoproto.stdtime) = true, (gogoproto.nullable) = false];
  Info info                                       = 11 [(gogoproto.nullable) = false];
}
```

### Lease
A lease represents an active lease of a media node by a user.
```protobuf
message Lease {
  string media_node_id                          = 1 [(gogoproto.moretags) = "yaml:\"media_node_id\""];
  string owner                                  = 2;
  string lessee                                 = 3;
  cosmos.base.v1beta1.Coin price_per_hour       = 4 [(gogoproto.nullable) = false, (gogoproto.moretags) = "yaml:\"price_per_hour\""];
  cosmos.base.v1beta1.Coin total_lease_amount   = 5 [(gogoproto.nullable) = false];
  cosmos.base.v1beta1.Coin settled_lease_amount = 6 [(gogoproto.moretags) = "yaml:\"settled_lease_amount\"", (gogoproto.nullable) = false];
  google.protobuf.Timestamp start_time          = 7 [(gogoproto.stdtime) = true, (gogoproto.nullable) = false, (gogoproto.moretags) = "yaml:\"start_time\""];
  uint64 leased_hours                           = 8 [(gogoproto.moretags) = "yaml:\"leased_hours\""];
  google.protobuf.Timestamp last_settled_at     = 9 [(gogoproto.stdtime) = true, (gogoproto.nullable) = false, (gogoproto.moretags) = "yaml:\"last_settled_at\""]; 
}
```

### Parameters
The module parameters include:
- Minimum deposit required for registration
- Initial deposit percentage
- Minimum and maximum lease hours
- Minimum price per hour
```protobuf
message Params {
  uint64 minimum_lease_hours = 1;
  uint64 maximum_lease_hours = 2;
  cosmos.base.v1beta1.Coin min_deposit       = 3 [
    (gogoproto.moretags)   = "yaml:\"min_deposit\"",
    (gogoproto.nullable)   = false
  ];
  string initial_deposit_percentage          = 4 [
    (gogoproto.moretags)   = "yaml:\"initial_deposit_percentage\"",
    (gogoproto.customtype) = "cosmossdk.io/math.LegacyDec",
    (gogoproto.nullable)   = false
  ];
  string lease_commission                    = 5 [
    (gogoproto.moretags)   = "yaml:\"lease_commission\"",
    (gogoproto.customtype) = "cosmossdk.io/math.LegacyDec",
    (gogoproto.nullable)   = false
  ];
  Distribution commission_distribution        = 6 [
    (gogoproto.moretags)   = "yaml:\"commission_distribution\"",
     (gogoproto.nullable) = false
  ];
  google.protobuf.Duration deposit_release_period = 7 [
    (gogoproto.nullable) = false,
    (gogoproto.stdduration) = true,
    (gogoproto.moretags) = "yaml:\"deposit_release_period\""
  ];
}
message Distribution {
  string staking        = 1 [
    (gogoproto.moretags)   = "yaml:\"staking\"",
    (gogoproto.customtype) = "cosmossdk.io/math.LegacyDec",
    (gogoproto.nullable)   = false
  ];
  string community_pool = 2 [
    (gogoproto.moretags)   = "yaml:\"community_pool\"",
    (gogoproto.customtype) = "cosmossdk.io/math.LegacyDec",
    (gogoproto.nullable)   = false
  ];
}
```

## Messages

`MsgRegisterMediaNode` can be used by any account to register a new media node.
`MsgUpdateMediaNode` can be used by the owner to update media node information.
`MsgLeaseMediaNode` can be used by any account to lease a media node.
`MsgExtendLease` can be used by the lessee to extend an existing lease.
`MsgCancelLease` can be used by the lessee to cancel an active lease.
`MsgDepositMediaNode` can be used by the owner to deposit additional funds to a media node.
`MsgCloseMediaNode` can be used by the owner to close a media node and withdraw deposits.

## Queries

The medianode module provides several queries to fetch information about media nodes and leases:

```go
service Query {
  rpc Params(QueryParamsRequest) returns (QueryParamsResponse) {
    option (google.api.http).get = "/omniflix/medianode/v1/params";
  }

  rpc MediaNode(QueryMediaNodeRequest) returns (QueryMediaNodeResponse) {
    option (google.api.http).get = "/omniflix/medianode/v1/nodes/{id}";
  }

  rpc MediaNodes(QueryMediaNodesRequest) returns (QueryMediaNodesResponse) {
    option (google.api.http).get = "/omniflix/medianode/v1/nodes";
  }

  rpc MediaNodesByOwner(QueryMediaNodesByOwnerRequest) returns (QueryMediaNodesByOwnerResponse) {
    option (google.api.http).get = "/omniflix/medianode/v1/nodes-by-owner/{owner}";
  }
  rpc Lease(QueryLeaseRequest) returns (QueryLeaseResponse) {
    option (google.api.http).get = "/omniflix/medianode/v1/lease/{media_node_id}";
  }

  rpc LeasesByLessee(QueryLeasesByLesseeRequest) returns (QueryLeasesByLesseeResponse) {
    option (google.api.http).get = "/omniflix/medianode/v1/leases-by-lessee/{lessee}";
  }
}
```