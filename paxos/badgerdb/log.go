package badgerdb

import (
	"encoding/binary"
	"reflect"

	"github.com/dgraph-io/badger/v3"
	"github.com/mjpitz/myago/paxos"
	"github.com/vmihailenco/msgpack/v5"
)

type Log struct {
	DB     *badger.DB
	Prefix []byte
}

func (l *Log) key(id uint64) []byte {
	prefixLen := len(l.Prefix)
	key := make([]byte, prefixLen+8)

	copy(key[:prefixLen], l.Prefix)
	binary.BigEndian.PutUint64(key[prefixLen:], id)

	return key
}

func (l *Log) Record(id uint64, msg interface{}) error {
	key := l.key(id)
	val, err := msgpack.Marshal(msg)
	if err != nil {
		return err
	}

	return l.DB.Update(func(txn *badger.Txn) error {
		return txn.Set(key, val)
	})
}

func (l *Log) Last(msg interface{}) error {
	txn := l.DB.NewTransaction(false)
	defer txn.Discard()

	iter := txn.NewIterator(badger.IteratorOptions{
		PrefetchSize:   1,
		PrefetchValues: true,
		Reverse:        true,
		Prefix:         l.Prefix,
	})
	defer iter.Close()

	if !iter.ValidForPrefix(l.Prefix) {
		return nil
	}

	return iter.Item().Value(func(val []byte) error {
		return msgpack.Unmarshal(val, msg)
	})
}

func (l *Log) Range(start, stop uint64, proto interface{}, fn func(msg interface{}) error) error {
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
	for iter.ValidForPrefix(l.Prefix) {
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

func (l *Log) Close() error {
	return l.DB.Close()
}

var _ paxos.Log = &Log{}
