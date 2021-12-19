# myago

My collection of Golang (Go) utilities for building distributed systems. Over the last few months, I've started to
consolidate some common code across my repositories into this single collection. It's allowed me to delete _some_
smaller repositories I have in favor of a single one with more of the common code I like to work with.

This isn't really intended to be used by others, but feel free explore, try things out, or submit issues if you find
them.

- [`authors`](authors) contains code for parsing [`AUTHORS`](https://opensource.google/docs/releasing/authors/) file contents.
- [`clocks`](clocks) contains code for working with clocks on contexts.
- [`cluster`](cluster) contains code for forming pools of nodes into clusters.
- [`config`](config) contains code for working with a variety of configuration file formats.
- [`dirset`](dirset) contains code for obtaining platform based, application state directories to cache, store, or log data to.
- [`encoding`](encoding) contains common encoding schemes used by configuration and for transport.
- [`flagset`](flagset) contains opinionated code for parsing Go structs into `urfave/cli` flags.
- [`leaderless`](leaderless) forms a `farm128` consistent hash ring to coordinate work within a cluster without the need for election.
- [`lifecycle`](lifecycle) provides hooks into the lifecycle of an application.
- [`livetls`](livetls) provides a `tls.Config` that periodically reloads the underlying configuration.
- [`paxos`](paxos) provides a paxos implementation.
- [`plugin`](plugin) provides code for writing command-line based plugins.
- [`ulid`](ulid) provides code for generating variable length unique, lexigraphic identifiers (ULID) with programmable fills.
- [`vfs`](vfs) provides code for working with file systems on contexts using afero.
- [`vue`](vue) contains helpers for VueJS applications.
- [`yarpc`](yarpc) is yet another RPC framework, built on top of Hashicorp Yamux with the simplicity of `http`.
- [`zaputil`](zaputil) is a collection of logging utilities for zap.

## Tooling

- (required) [`go`](https://golang.org/). See `go.mod` for required version information.
- (required) [`openssl`](https://www.openssl.org/) is used to generate certificates for tests in the `livetls` package.
- (optional) [`addlicense`](https://github.com/google/addlicense) prepends files with appropriate license information.
