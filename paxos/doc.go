/*
Package paxos implements the paxos algorithm. The logic is mostly ported from mjpitz/paxos, but with a few
modifications. First, I didn't continue using gRPC as the transport as I wanted something a bit less cumbersome. I've
tried to break down the interface in such a way where different transports _could_ be plugged in.

This package is still a work in progress. The current code block supports a single acceptor process and requires
modifications (mostly behind existing interfaces) to support multiple acceptors.
*/
package paxos
