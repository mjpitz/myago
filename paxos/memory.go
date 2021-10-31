package paxos

import (
	"reflect"
	"sort"
	"sync"

	"github.com/vmihailenco/msgpack/v5"
)

type Memory struct {
	mu     sync.RWMutex
	idLog  []uint64
	msgLog [][]byte
}

func (m *Memory) WithPrefix(prefix string) Log {
	// return a new memory instance for the prefix
	return &Memory{}
}

func (m *Memory) Record(id uint64, msg interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, err := msgpack.Marshal(msg)
	if err != nil {
		return err
	}

	idx := sort.Search(len(m.idLog), func(i int) bool {
		return id <= m.idLog[i]
	})

	// in paxos, the last two cases of this switch statement should _never_ happen
	switch {
	case idx == len(m.idLog):
		m.idLog = append(m.idLog, id)
		m.msgLog = append(m.msgLog, data)
	case m.idLog[idx] == id:
		// entry exists in log
		return nil
	default:
		m.idLog = append(m.idLog[:idx], append([]uint64{id}, m.idLog[idx:]...)...)
		m.msgLog = append(m.msgLog[:idx], append([][]byte{data}, m.msgLog[idx:]...)...)
	}
	return nil
}

func (m *Memory) Last(msg interface{}) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.msgLog) == 0 {
		return nil
	}

	return msgpack.Unmarshal(m.msgLog[len(m.msgLog)-1], msg)
}

func (m *Memory) Range(start, end uint64, proto interface{}, fn func(msg interface{}) error) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	startIdx := sort.Search(len(m.idLog), func(i int) bool {
		return start <= m.idLog[i]
	})

	endIdx := sort.Search(len(m.idLog), func(i int) bool {
		return end <= m.idLog[i]
	})

	if startIdx == len(m.idLog) {
		return nil
	}

	for i := startIdx; i <= endIdx; i++ {
		inst := reflect.New(reflect.TypeOf(proto)).Interface()
		err := msgpack.Unmarshal(m.msgLog[i], inst)
		if err != nil {
			return err
		}

		err = fn(inst)
		if err != nil {
			return err
		}
	}

	return nil
}

var _ Log = &Memory{}
