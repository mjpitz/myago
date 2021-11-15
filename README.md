# myago

My collection of Golang (Go) utilities for building distributed systems. Not really intended to be used by others, but
feel free to poke around.

- [`authors`](authors) contains code for parsing `AUTHORS` file contents.
- [`clocks`](clocks) contains code for managing clocks on contexts.
- [`cluster`](cluster) contains code for forming pools of nodes into clusters.
- [`encoding`](encoding) contains structures for common encoding schemes.
- [`flagset`](flagset) contains opinionated code for parsing Go structs into urfave/cli flags.
- [`leaderless`](leaderless) forms a `farm128` consistent hash ring to coordinate work within a cluster without the need
  for election.
- [`lifecycle`](lifecycle) provides hooks into the lifecycle of an application.
- [`livetls`](livetls) provides code for creating a `tls.Config` that periodically reloads it's configuration.
- [`paxos`](paxos) provides a paxos implementation.
- [`ulid`](ulid) provides code for generating variable length unique, lexigraphic identifiers (ULID) with programmable
  fills.
- [`vfs`](vfs) provides code for managing file systems through afero.
- [`yarpc`](yarpc) is yet another RPC framework, built on top of HashiCorp Yamux with the simplicity of `http`.
- [`zaputil`](zaputil) is a collection of logging utilities for zap.

## Required Tooling

- [`go`](https://golang.org/). See `go.mod` for required version information.
- [`openssl`](https://www.openssl.org/) is used to generation certificates for the `livetls` package.
