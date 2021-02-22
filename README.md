[![GitHub](https://img.shields.io/github/license/hashgraph/hedera-state-proof-verifier-go)](LICENSE)
[![Discord](https://img.shields.io/badge/discord-join%20chat-blue.svg)](https://hedera.com/discord)

# Hedera Verify State Proof in Go

Cryptographically prove a transaction is valid on Hedera Network in Go.
Based on [official documentation](https://docs.hedera.com/guides/docs/record-and-event-stream-file-formats).

# Install

```
go get https://github.com/hashgraph/hedera-state-proof-verifier-go
```

# How to use?

```go
import "github.com/hashgraph/hedera-state-proof-verifier-go/stateproof"

verified, err := stateproof.Verify(txnID, stateProof)
```

# Examples

[V2 Record Stream State Proof Verification](examples/v2/main.go)

[V5 Record Stream State Proof Verification](examples/v5/main.go)
