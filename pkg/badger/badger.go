package badger

import (
	"sync/atomic"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Store is a wrapper around a badger DB
type Store struct {
	db *badger.DB
}

var stor *Store

func Open(dbPath string, logger *zap.Logger) (*Store, error) {
	opts := badger.DefaultOptions("")

	opts.Logger = NewLogger(logger)

	if dbPath == ":memory:" || dbPath == "" {
		// In-Memory Mode/Diskless Mode
		opts.InMemory = true
	} else {
		// Disk Mode
		opts.Dir = dbPath
		opts.ValueDir = dbPath
	}
	db, err := badger.Open(opts) //文件只能被一个进程使用，如果不调用Close则下次无法Open。手动释放锁的办法：把LOCK文件删掉
	if err != nil {
		panic(err)
	}
	stor = &Store{db}
	if err != nil {
		panic(err)
	}
	return stor, err
}

func (s *Store) CreateBucket(name string) *Bucket {
	return &Bucket{
		name:   name,
		prefix: name + ":",
		stor:   stor,
	}
}

func (s *Store) Set(k, v []byte) error {
	err := s.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(k, v)
		return err
	})
	return err
}

func (s *Store) SetWithTTL(k, v []byte, expireAt int64) error {
	err := s.db.Update(func(txn *badger.Txn) error {
		duration := time.Duration(expireAt-time.Now().Unix()) * time.Second
		e := badger.NewEntry(k, v).WithTTL(duration)
		err := txn.SetEntry(e)
		return err
	})
	return err
}

//BatchSet 多个写操作使用一个事务
func (s *Store) BatchSet(keys, values [][]byte) error {
	if len(keys) != len(values) {
		return errors.New("key value not the same length")
	}

	// Start a writable transaction.
	txn := s.db.NewTransaction(true)
	defer txn.Discard()

	for i, key := range keys {
		value := values[i]
		// Use the transaction...
		err := txn.Set(key, value)
		return err
	}
	// Commit the transaction and check for error.
	err := txn.Commit()
	return err
}

//BatchSet 多个写操作使用一个事务
func (s *Store) BatchSetWithTTL(keys, values [][]byte, expireAts []int64) error {
	if len(keys) != len(values) {
		return errors.New("key value not the same length")
	}
	// Start a writable transaction.
	txn := s.db.NewTransaction(true)
	defer txn.Discard()
	for i, key := range keys {
		value := values[i]
		duration := time.Duration(expireAts[i]-time.Now().Unix()) * time.Second
		// Use the transaction...
		e := badger.NewEntry(key, value).WithTTL(duration)
		err := txn.SetEntry(e)
		return err
	}
	// Commit the transaction and check for error.
	err := txn.Commit()
	return err
}

//Get 如果key不存在会返回error:Key not found
func (s *Store) Get(k []byte) ([]byte, error) {
	var ival []byte
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(k)
		if err != nil {
			return err
		}
		ival, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}
		return nil
	})
	return ival, err
}

//BatchGet 返回的values与传入的keys顺序保持一致。如果key不存在或读取失败则对应的value是空数组
func (s *Store) BatchGet(keys [][]byte) ([][]byte, error) {
	var err error
	txn := s.db.NewTransaction(false) //只读事务
	defer txn.Discard()
	values := make([][]byte, len(keys))
	for i, key := range keys {
		var item *badger.Item
		item, err = txn.Get(key)
		if err == nil {
			var ival []byte
			ival, err = item.ValueCopy(nil)
			if err == nil {
				values[i] = ival
			} else { //拷贝失败
				values[i] = []byte{} //拷贝失败就把value设为空数组
			}
		} else { //读取失败
			values[i] = []byte{}              //读取失败就把value设为空数组
			if err != badger.ErrKeyNotFound { //如果真的发生异常，则开一个新事务继续读后面的key
				txn.Discard()
				txn = s.db.NewTransaction(false)
			}
		}
	}
	return values, err
}

//Delete
func (s *Store) Delete(k []byte) error {
	err := s.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(k)
	})
	return err
}

//BatchDelete
func (s *Store) BatchDelete(keys [][]byte) error {
	var err error
	txn := s.db.NewTransaction(true)
	defer txn.Discard()
	for _, key := range keys {
		if err = txn.Delete(key); err != nil {
			_ = txn.Commit() //发生异常时就提交老事务，然后开一个新事务，重试delete
			txn = s.db.NewTransaction(true)
			_ = txn.Delete(key)
		}
	}
	_ = txn.Commit()
	return err
}

//Has 判断某个key是否存在
func (s *Store) Has(k []byte) bool {
	var exists = false
	s.db.View(func(txn *badger.Txn) error { //Store.View相当于打开了一个读写事务:Store.NewTransaction(true)。用db.Update的好处在于不用显式调用Txn.Discard()了
		_, err := txn.Get(k)
		if err != nil {
			return err
		} else {
			exists = true //没有任何异常发生，则认为k存在。如果k不存在会发生ErrKeyNotFound
		}
		return err
	})
	return exists
}

//IterDB 遍历整个DB
func (s *Store) Iter(fn func(k, v []byte) error) int64 {
	var total int64
	s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := item.Key()
			val, err := item.ValueCopy(nil)
			if err != nil {
				continue
			}
			if err := fn(key, val); err == nil {
				atomic.AddInt64(&total, 1)
			}
		}
		return nil
	})
	return atomic.LoadInt64(&total)
}

func (s *Store) IterByPrefix(prefix string, fn func(k, v []byte) error) int64 {
	var total int64
	s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		prefix := []byte(prefix)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			key := item.Key()
			val, err := item.ValueCopy(nil)
			if err != nil {
				continue
			}
			if err := fn(key, val); err == nil {
				atomic.AddInt64(&total, 1)
			}
		}
		return nil
	})
	return atomic.LoadInt64(&total)
}

//IterKey 只遍历key。key是全部存在LSM tree上的，只需要读内存，所以很快
func (s *Store) IterKeys(fn func(k []byte) error) int64 {
	var total int64
	s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false //只需要读key，所以把PrefetchValues设为false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			if err := fn(k); err == nil {
				atomic.AddInt64(&total, 1)
			}
		}
		return nil
	})
	return atomic.LoadInt64(&total)
}

func (s *Store) IterKeysByPrefix(prefix string, fn func(k []byte) error) int64 {
	var total int64
	s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false //只需要读key，所以把PrefetchValues设为false
		it := txn.NewIterator(opts)
		defer it.Close()
		prefix := []byte(prefix)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()
			if err := fn(k); err == nil {
				atomic.AddInt64(&total, 1)
			}
		}
		return nil
	})
	return atomic.LoadInt64(&total)
}

func (s *Store) CheckAndGC() {
	for {
		if err := s.db.RunValueLogGC(0.5); err == badger.ErrNoRewrite || err == badger.ErrRejected {
			break
		}
	}
}

func (s *Store) Size() (int64, int64) {
	return s.db.Size()
}

//Clean 清空所有数据
func (s *Store) Clean() error {
	return s.db.DropAll()
}

//Close 把内存中的数据flush到磁盘，同时释放文件锁
func (s *Store) Close() error {
	return s.db.Close()
}
