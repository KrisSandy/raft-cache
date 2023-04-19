package cache

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	bolt "go.etcd.io/bbolt"
)

type store interface {
	Get(bucket string, key string) (string, bool)
	Put(bucket string, key string, value string) error
	CreateBucket(name string) error
	Close() error
	GetDBPath() string
	Restore(rc io.ReadCloser) error
}

type storeDB struct {
	db   *bolt.DB
	path string

	mtx sync.RWMutex
}

func newStore(path string) (store, error) {
	db, err := bolt.Open(filepath.Join(path, DBName), 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open cache db: %v", err)
	}

	return &storeDB{
		db:   db,
		path: path,
	}, nil
}

func (s *storeDB) Get(bucket string, key string) (string, bool) {
	var v []byte

	s.mtx.RLock()
	defer s.mtx.RUnlock()

	s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		v = b.Get([]byte(key))
		return nil
	})

	if v == nil {
		return "", false
	}

	return string(v), true
}

func (s *storeDB) Put(bucket string, key string, value string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Put([]byte(key), []byte(value))
		return err
	})

	return nil
}

func (s *storeDB) CreateBucket(name string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(name))
		if err != nil {
			return err
		}
		return nil
	})
	return nil
}

func (s *storeDB) Close() error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	return s.db.Close()
}

func (s *storeDB) GetDBPath() string {
	return s.path
}

func (s *storeDB) Restore(rc io.ReadCloser) error {
	bytes, err := io.ReadAll(rc)
	if err != nil {
		return fmt.Errorf("failed to read snapshot: %v", err)
	}

	if err = os.WriteFile(filepath.Join(s.GetDBPath(), DBName), bytes, 0644); err != nil {
		return fmt.Errorf("failed to write snapshot: %v", err)
	}

	return nil
}
