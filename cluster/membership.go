package cluster

import (
	"sort"
	"sync"
)

// Membership tacks a current list of active within the cluster. It can be populated manually (useful for testing) or
// using common discovery mechanisms.
type Membership struct {
	mu      sync.Mutex
	active  []string
	left    []string
	watches map[*struct{}]chan MembershipChange
}

// broadcast requires external locking of the sync.Mutex before attempting to send a notification. Otherwise, you may
// encounter a case where the underlying map is concurrently modified (likely by an unsubscription).
func (m *Membership) broadcast(notification MembershipChange) {
	for _, ch := range m.watches {
		ch <- notification
	}
}

// Add inserts the provided active into the cluster's active list.
// Operation should be `O( m log(n) )` where `m = len(peers)` and `n = len(m.active) + len(m.left)`.
func (m *Membership) Add(peers []string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	membershipChange := MembershipChange{Active: peers}
	defer m.broadcast(membershipChange)

	for _, peer := range peers {
		// remove from left
		leftIdx := sort.SearchStrings(m.left, peer)
		switch {
		case leftIdx == len(m.left):
			// pass
		case m.left[leftIdx] == peer:
			m.left = append(m.left[:leftIdx], m.left[leftIdx+1:]...)
		}

		// add to active
		activeIdx := sort.SearchStrings(m.active, peer)
		switch {
		case activeIdx == len(m.active):
			m.active = append(m.active, peer)
		case activeIdx == 0:
			m.active = append([]string{peer}, m.active...)
		case m.active[activeIdx] != peer:
			m.active = append(m.active[:activeIdx], append([]string{peer}, m.active[activeIdx:]...)...)
		default:
			// otherwise already exists ...
		}
	}
}

// Left allows peers to temporarily leave the cluster, but still be considered part of active membership.
// Operation should be `O( m log(n) )` where `m = len(peers)` and `n = len(m.active) + len(m.left)`.
func (m *Membership) Left(peers []string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	membershipChange := MembershipChange{Left: peers}
	defer m.broadcast(membershipChange)

	for _, peer := range peers {
		// remove from active
		activeIdx := sort.SearchStrings(m.active, peer)
		switch {
		case activeIdx == len(m.active):
			continue
		case m.active[activeIdx] != peer:
			continue
		case m.active[activeIdx] == peer:
			m.active = append(m.active[:activeIdx], m.active[activeIdx+1:]...)
		}

		// add to left
		leftIdx := sort.SearchStrings(m.left, peer)
		switch {
		case leftIdx == len(m.left):
			m.left = append(m.left, peer)
		case leftIdx == 0:
			m.left = append([]string{peer}, m.left...)
		case m.left[leftIdx] != peer:
			m.left = append(m.left[:leftIdx], append([]string{peer}, m.left[leftIdx:]...)...)
		default:
			// otherwise already exists ...
		}
	}
}

// Remove deletes the provided active from the cluster's peer list.
// Operation should be `O( m log(n) )` where `m = len(peers)` and `n = len(m.active) + len(m.left)`.
func (m *Membership) Remove(peers []string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	membershipChange := MembershipChange{Removed: peers}
	defer m.broadcast(membershipChange)

	for _, peer := range peers {
		// remove from active
		activeIdx := sort.SearchStrings(m.active, peer)
		switch {
		case activeIdx == len(m.active):
			// pass
		case m.active[activeIdx] == peer:
			m.active = append(m.active[:activeIdx], m.active[activeIdx+1:]...)
		}

		// remove from left
		leftIdx := sort.SearchStrings(m.left, peer)
		switch {
		case leftIdx == len(m.left):
			// pass
		case m.left[leftIdx] == peer:
			m.left = append(m.left[:leftIdx], m.left[leftIdx+1:]...)
		}
	}
}

func (m *Membership) Majority() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return (len(m.active) + len(m.left) + 1) / 2
}

// Watch allows others to observe changes in the cluster membership.
func (m *Membership) Watch() (<-chan MembershipChange, CancelWatch) {
	id := &struct{}{}
	watchChan := make(chan MembershipChange, 3)

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.watches == nil {
		m.watches = make(map[*struct{}]chan MembershipChange)
	}
	m.watches[id] = watchChan

	watchChan <- MembershipChange{
		Active: m.active,
		Left:   m.left,
	}

	return watchChan, func() {
		m.mu.Lock()
		defer m.mu.Unlock()

		delete(m.watches, id)
	}
}

// Snapshot returns a copy of the current peer list.
func (m *Membership) Snapshot() ([]string, int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	la := len(m.active)
	ll := len(m.left)

	peers := make([]string, la+ll)

	split := copy(peers[:la], m.active)
	n := split + copy(peers[split:split+ll], m.left)

	return peers[:n], split
}

// MembershipChange describes how the cluster membership has changed to outside observers.
type MembershipChange struct {
	Active  []string
	Left    []string
	Removed []string
}

// CancelWatch is used to remove a watch from the cluster membership.
type CancelWatch func()
