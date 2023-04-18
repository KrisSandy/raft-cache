package cache

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/hashicorp/raft"
)

type cacheSnapshot struct {
	items map[string]string
}

var _ raft.FSMSnapshot = (*cacheSnapshot)(nil)

func (s *cacheSnapshot) Persist(sink raft.SnapshotSink) error {
	defer sink.Close()

	buf, err := GetBytes(s.items)
	if err != nil {
		sink.Cancel()
		return fmt.Errorf("failed to encode items: %v", err)
	}

	if _, err := sink.Write(buf); err != nil {
		sink.Cancel()
		return fmt.Errorf("failed to write snapshot: %v", err)
	}

	if err := sink.Close(); err != nil {
		return fmt.Errorf("failed to close snapshot: %v", err)
	}

	return nil
}

func (s *cacheSnapshot) Release() {}

func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
