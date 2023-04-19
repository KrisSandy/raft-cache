package cache

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"

	"github.com/hashicorp/raft"
)

const (
	OpSet Operation = iota
	OpCreateBucket
)

type fsm struct {
	store store
}

type Operation int

type command struct {
	Op     Operation
	Bucket string
	Key    string
	Value  string
}

var _ raft.FSM = (*fsm)(nil)

func newFSM(path string) (*fsm, error) {
	store, err := newStore(path)
	if err != nil {
		return nil, err
	}

	return &fsm{
		store: store,
	}, nil
}

func (c *fsm) Get(bucket string, key string) (string, bool) {
	value, ok := c.store.Get(bucket, key)
	return value, ok
}

func (c *fsm) Apply(l *raft.Log) interface{} {
	data := bytes.NewBuffer(l.Data)
	var cmd command
	if err := gob.NewDecoder(data).Decode(&cmd); err != nil {
		log.Println("failed to decode command")
		return fmt.Errorf("failed to decode command: %v", err)
	}

	log.Println("Applying command to cache: %w", cmd)

	switch cmd.Op {
	case OpSet:
		err := c.store.Put(cmd.Bucket, cmd.Key, cmd.Value)
		if err != nil {
			return fmt.Errorf("failed to set %q: %v", cmd.Key, err)
		}
	case OpCreateBucket:
		err := c.store.CreateBucket(cmd.Bucket)
		if err != nil {
			return fmt.Errorf("failed to create bucket %q: %v", cmd.Bucket, err)
		}
	default:
		return fmt.Errorf("unknown command: %v", cmd.Op)
	}

	return nil
}

func (c *fsm) Snapshot() (raft.FSMSnapshot, error) {
	return &cacheSnapshot{
		store: c.store,
	}, nil
}

func (c *fsm) Restore(rc io.ReadCloser) error {
	return c.store.Restore(rc)
}

func (c *fsm) Close() error {
	return c.store.Close()
}
