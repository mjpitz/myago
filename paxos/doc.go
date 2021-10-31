/*
Package paxos implements the paxos algorithm. The logic is mostly ported from mjpitz/paxos, but with a few
modifications. First, I didn't continue using gRPC as the transport as I wanted something a bit less cumbersome. I've
tried to break down the interface in such a way where different transports _could_ be plugged in. More on that later.

This package is (and likely will be for a while) a work in progress. As it stands, it _should_ support simple paxos.
*/
package paxos
