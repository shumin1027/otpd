package badger

import (
	"context"
	"github.com/dgraph-io/badger/v3"
	"github.com/dgraph-io/badger/v3/pb"
	"github.com/dgraph-io/ristretto/z"
	"log"
	"runtime"
	"sync/atomic"
)

type Bucket struct {
	name   string
	prefix string
	stor   *Store
}

func CreateBucket(name string, stor *Store) *Bucket {
	return &Bucket{
		name:   name,
		prefix: name + ":",
		stor:   stor,
	}
}

func (s *Bucket) Set(k, v []byte) error {
	k = []byte(s.prefix + string(k))
	return s.stor.Set(k, v)
}

func (s *Bucket) SetWithTTL(k, v []byte, expireAt int64) error {
	k = []byte(s.prefix + string(k))
	return s.stor.SetWithTTL(k, v, expireAt)
}

//BatchSet 多个写操作使用一个事务
func (s *Bucket) BatchSet(keys, values [][]byte) error {
	for i := 0; i < len(keys); i++ {
		k := keys[i]
		k = []byte(s.prefix + string(k))
		keys[i] = k
	}
	return s.stor.BatchSet(keys, values)
}

//BatchSet 多个写操作使用一个事务
func (s *Bucket) BatchSetWithTTL(keys, values [][]byte, expireAts []int64) error {
	for i := 0; i < len(keys); i++ {
		k := keys[i]
		k = []byte(s.prefix + string(k))
		keys[i] = k
	}
	return s.stor.BatchSetWithTTL(keys, values, expireAts)
}

//Get 如果key不存在会返回error:Key not found
func (s *Bucket) Get(k []byte) ([]byte, error) {
	k = []byte(s.prefix + string(k))
	return s.stor.Get(k)
}

//BatchGet 返回的values与传入的keys顺序保持一致。如果key不存在或读取失败则对应的value是空数组
func (s *Bucket) BatchGet(keys [][]byte) ([][]byte, error) {
	for i := 0; i < len(keys); i++ {
		k := keys[i]
		k = []byte(s.prefix + string(k))
		keys[i] = k
	}
	return s.stor.BatchGet(keys)
}

//Delete
func (s *Bucket) Delete(k []byte) error {
	k = []byte(s.prefix + string(k))
	return s.stor.Delete(k)
}

//BatchDelete
func (s *Bucket) BatchDelete(keys [][]byte) error {
	for i := 0; i < len(keys); i++ {
		k := keys[i]
		k = []byte(s.prefix + string(k))
		keys[i] = k
	}
	return s.stor.BatchDelete(keys)
}

//Has 判断某个key是否存在
func (s *Bucket) Has(k []byte) bool {
	k = []byte(s.prefix + string(k))
	return s.stor.Has(k)
}

//IterDB 遍历整个DB
func (s *Bucket) Iter(fn func(k, v []byte) error) int64 {
	return s.stor.IterByPrefix(s.prefix, fn)
}

//IterKey 只遍历key。key是全部存在LSM tree上的，只需要读内存，所以很快
func (s *Bucket) IterKeys(fn func(k []byte) error) int64 {
	return s.stor.IterKeysByPrefix(s.prefix, fn)
}

func (s *Bucket) Stream(fn func(k []byte, v []byte) error) int64 {
	var total int64
	stream := s.stor.db.NewStream()
	// -- Optional settings
	stream.NumGo = runtime.NumCPU() * 2   // Set number of goroutines to use for iteration.
	stream.Prefix = []byte(s.prefix)      // Leave nil for iteration over the whole DB.
	stream.LogPrefix = "Badger.Streaming" // For identifying stream logs. Outputs to Logger.

	stream.KeyToList = func(key []byte, itr *badger.Iterator) (*pb.KVList, error) {
		item := itr.Item()
		val, err := item.ValueCopy(nil)
		if err != nil {
			return nil, err
		}
		err = fn(key, val)
		if err != nil {
			return nil, err
		}
		atomic.AddInt64(&total, 1)
		return stream.ToList(key, itr)
	}
	stream.Send = func(buf *z.Buffer) error {
		return nil
	}
	// Run the stream
	if err := stream.Orchestrate(context.Background()); err != nil {
		log.Panic(err)
	}
	// Done.
	return atomic.LoadInt64(&total)
}

func (s *Bucket) Chan(ch chan<- interface{}) int64 {
	var total int64
	stream := s.stor.db.NewStream()
	// -- Optional settings
	stream.NumGo = runtime.NumCPU() * 2   // Set number of goroutines to use for iteration.
	stream.Prefix = []byte(s.prefix)      // Leave nil for iteration over the whole DB.
	stream.LogPrefix = "Badger.Streaming" // For identifying stream logs. Outputs to Logger.

	stream.KeyToList = func(key []byte, itr *badger.Iterator) (*pb.KVList, error) {
		item := itr.Item()
		val, err := item.ValueCopy(nil)
		if err != nil {
			return nil, err
		}

		// send value to channel
		ch <- val

		atomic.AddInt64(&total, 1)
		return stream.ToList(key, itr)
	}
	stream.Send = func(buf *z.Buffer) error {
		return nil
	}
	// Run the stream
	if err := stream.Orchestrate(context.Background()); err != nil {
		log.Panic(err)
	}
	// Done.
	return atomic.LoadInt64(&total)
}
