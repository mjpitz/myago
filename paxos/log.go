package paxos

type Log interface {
	WithPrefix(str string) Log
	Record(id uint64, msg interface{}) error
	Last(msg interface{}) error
	Range(start, stop uint64, proto interface{}, fn func(msg interface{}) error) error
}
