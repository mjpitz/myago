package paxos

import (
	"reflect"
	"sort"
	"sync"

	"github.com/vmihailenco/msgpack/v5"
)

type Log interface {
	Record(id uint64, msg interface{}) error
	Last(msg interface{}) error
	Range(start, stop uint64, proto interface{}, fn func(msg interface{}) error) error
	Close() error
}

type MemoryLog struct {
	mu     sync.RWMutex
	idLog  []uint64
	msgLog [][]byte
}

func (m *MemoryLog) Record(id uint64, msg interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, err := msgpack.Marshal(msg)
	if err != nil {
		return err
	}

	m.idLog = append(m.idLog, id)
	m.msgLog = append(m.msgLog, data)
	return nil
}

func (m *MemoryLog) Last(msg interface{}) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.msgLog) == 0 {
		return nil
	}

	return msgpack.Unmarshal(m.msgLog[len(m.msgLog)-1], msg)
}

func (m *MemoryLog) Range(start, end uint64, proto interface{}, fn func(msg interface{}) error) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	startIdx := sort.Search(len(m.idLog), func(i int) bool {
		return start <= m.idLog[i]
	})

	endIdx := sort.Search(len(m.idLog), func(i int) bool {
		return end <= m.idLog[i]
	})

	for i := startIdx ; i <= endIdx ; i++ {
		inst := reflect.New(reflect.TypeOf(proto)).Interface()
		err := msgpack.Unmarshal(m.msgLog[i], inst)
		if err != nil {
			return err
		}

		err = fn(inst)
	}

	return nil
}

func (m *MemoryLog) Close() error {
	return nil
}

var _ Log = &MemoryLog{}
