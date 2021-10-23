# myago

My collection of Golang (Go) utilities for building distributed systems

- `cluster` contains code for forming pools of nodes into clusters.
- `leaderless` forms a `farm128` consistent hash ring to coordinate work within a cluster without the need for election.
- `ulid256` provides code for generating 256bit unique, lexigraphic identifiers (ULID).
- `yarpc` is yet another RPC framework, built on top of HashiCorp Yamux with the simplicity of `http`.
