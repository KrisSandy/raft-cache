package cache

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/hashicorp/raft"
	"golang.org/x/exp/maps"
)

const (
	OpSet Operation = iota
	OpDelete
)

type cache struct {
	mtx   sync.RWMutex
	items map[string]string
}

type Operation int

type command struct {
	Op    Operation
	Key   string
	Value string
}

var _ raft.FSM = (*cache)(nil)

func newCache() *cache {
	return &cache{
		items: make(map[string]string),
	}
}

func (c *cache) Get(key string) (string, bool) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	value, ok := c.items[key]
	return value, ok
}

func (c *cache) Apply(l *raft.Log) interface{} {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	data := bytes.NewBuffer(l.Data)
	var cmd command
	if err := gob.NewDecoder(data).Decode(&cmd); err != nil {
		log.Println("failed to decode command")
		return fmt.Errorf("failed to decode command: %v", err)
	}

	log.Println("Applying command to cache: %w", cmd)

	switch cmd.Op {
	case OpSet:
		log.Printf("setting %q to %q", cmd.Key, cmd.Value)
		c.items[cmd.Key] = cmd.Value
	case OpDelete:
		log.Printf("deleting %q", cmd.Key)
		delete(c.items, cmd.Key)
	default:
		return fmt.Errorf("unknown command: %v", cmd.Op)
	}

	return nil
}

func (c *cache) Snapshot() (raft.FSMSnapshot, error) {
	return &cacheSnapshot{
		items: copyMap(c.items),
	}, nil
}

func (c *cache) Restore(rc io.ReadCloser) error {
	dec := gob.NewDecoder(rc)
	defer rc.Close()

	var items map[string]string
	if err := dec.Decode(&items); err != nil {
		return fmt.Errorf("failed to decode items: %v", err)
	}

	c.items = items
	return nil
}

func copyMap(m map[string]string) map[string]string {
	n := make(map[string]string, len(m))
	maps.Copy(n, m)
	return n
}
