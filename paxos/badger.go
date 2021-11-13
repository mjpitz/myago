package paxos

import (
	"encoding/binary"
	"errors"
	"reflect"

	"github.com/dgraph-io/badger/v3"
	"github.com/vmihailenco/msgpack/v5"
)

// Badger implements a Log that wraps an underlying badgerdb instance.
type Badger struct {
	DB     *badger.DB
	prefix []byte
}

func (l *Badger) key(id uint64) []byte {
	prefixLen := len(l.prefix)
	key := make([]byte, prefixLen+8)

	copy(key[:prefixLen], l.prefix)
	binary.BigEndian.PutUint64(key[prefixLen:], id)

	return key
}

func (l *Badger) WithPrefix(prefix string) Log {
	log := &Badger{
		DB: l.DB,
	}

	log.prefix = append(log.prefix, l.prefix...)
	log.prefix = append(log.prefix, []byte(prefix)...)

	return log
}

func (l *Badger) Record(id uint64, msg interface{}) error {
	key := l.key(id)
	val, err := msgpack.Marshal(msg)
	if err != nil {
		return err
	}

	return l.DB.Update(func(txn *badger.Txn) error {
		_, err := txn.Get(key)
		if errors.Is(err, badger.ErrKeyNotFound) {
			err = txn.Set(key, val)
		}

		return err
	})
}

func (l *Badger) Last(msg interface{}) error {
	txn := l.DB.NewTransaction(false)
	defer txn.Discard()

	// find a better way to do this...

	iter := txn.NewIterator(badger.IteratorOptions{
		PrefetchSize: 1,
		Reverse:      true,
	})
	defer iter.Close()

	prefix := append([]byte{}, l.prefix...)
	prefix = append(prefix, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff)

	iter.Seek(prefix)
	if !iter.ValidForPrefix(l.prefix) {
		return nil
	}

	return iter.Item().Value(func(val []byte) error {
		return msgpack.Unmarshal(val, msg)
	})
}

func (l *Badger) Range(start, stop uint64, proto interface{}, fn func(msg interface{}) error) error {
	startKey := l.key(start)
	stopKey := l.key(stop)

	txn := l.DB.NewTransaction(false)
	defer txn.Discard()

	iter := txn.NewIterator(badger.IteratorOptions{
		PrefetchSize:   int(stop - start),
		PrefetchValues: true,
	})
	defer iter.Close()

	iter.Seek(startKey)
	for iter.ValidForPrefix(l.prefix) {
		inst := reflect.New(reflect.TypeOf(proto)).Interface()

		err := iter.Item().Value(func(val []byte) error {
			return msgpack.Unmarshal(val, inst)
		})
		if err != nil {
			return err
		}

		err = fn(inst)
		if err != nil {
			return err
		}

		if iter.ValidForPrefix(stopKey) {
			break
		}
		iter.Next()
	}

	return nil
}

var _ Log = &Badger{}
