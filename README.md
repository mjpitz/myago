# myago

My collection of Golang (Go) utilities for building distributed systems. Not really intended to be used by others, but feel free to poke around.

- [`authors`][] contains code for parsing `AUTHORS` file contents.
- [`clocks`][] contains code for managing clocks on contexts.
- [`cluster`][] contains code for forming pools of nodes into clusters.
- [`flagset`][] contains opinionated code for parsing Go structs into urfave/cli flags.
- [`leaderless`][] forms a `farm128` consistent hash ring to coordinate work within a cluster without the need for election.
- [`lifecycle`][] provides hooks into the lifecycle of an application.
- [`paxos`][] provides a paxos implementation.
- [`ulid`][] provides code for generating variable length unique, lexigraphic identifiers (ULID) with programmable fills.
- [`vfs`][] provides code for managing file systems through afero.
- [`yarpc`][] is yet another RPC framework, built on top of HashiCorp Yamux with the simplicity of `http`.
- [`zaputil`][] is a collection of logging utilities for zap.

[`authors`]: authors
[`clocks`]: clocks
[`cluster`]: cluster
[`flagset`]: flagset
[`leaderless`]: leaderless
[`lifecycle`]: lifecycle
[`paxos`]: paxos
[`ulid`]: ulid
[`vfs`]: vfs
[`yarpc`]: yarpc
[`zaputil`]: zaputil
