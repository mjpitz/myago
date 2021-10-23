/*
Package leaderless implements leader election without the need for coordination. It does this using a highly stable set
of peers and a hash ring. Since we know that each node shares the same view of the world, then we know that the computed
hash ring will be the same between instances. However, if peer knowledge is not stable (for example, members come and go
freely) then leaderless can result in a split brain state where some nodes share a different view of the world until the
other nodes "catch up".

This package is loosely inspired by Uber's ringpop system which seems like it's used quite extensively.
*/
package leaderless
