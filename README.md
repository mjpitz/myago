# myago

My collection of Golang (Go) utilities for building distributed systems. Not really intended to be used by others, but feel free to poke around.

- [`cluster`][] contains code for forming pools of nodes into clusters.
- [`leaderless`][] forms a `farm128` consistent hash ring to coordinate work within a cluster without the need for election.
- [`ulid256`][] provides code for generating 256bit unique, lexigraphic identifiers (ULID).
- [`yarpc`][] is yet another RPC framework, built on top of HashiCorp Yamux with the simplicity of `http`.

[`cluster`]: cluster
[`leaderless`]: leaderless
[`ulid256`]: ulid256
[`yarpc`]: yarpc
