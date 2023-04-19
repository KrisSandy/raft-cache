package cache

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"example.com/raft-cache/proto/cs"

	"example.com/raft-cache/pkg/admin"

	transport "github.com/Jille/raft-grpc-transport"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb/v2"
	"google.golang.org/grpc"
)

type Cache interface {
	Get(bucket string, key string) (string, bool)
	Put(bucket string, key string, value string) error
	CreateBucket(bucket string) error
}

type cacheNode struct {
	r      *raft.Raft
	c      *fsm
	voters map[string]struct{}
}

func NewCacheNode(raftID string, raftAddr string, raftDataDir string, bootstrap bool, voters map[string]struct{}) (Cache, error) {
	c, err := newFSM(raftDataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create new FSM: %v", err)
	}

	if !bootstrap {
		bootstrap = checkForAutomaticBootstrap(raftID, raftDataDir)
	}

	r, tm, err := newRaftNode(raftID, raftAddr, raftDataDir, bootstrap, c)
	if err != nil {
		return nil, fmt.Errorf("failed to create new raft node: %v", err)
	}

	cn := &cacheNode{
		r:      r,
		c:      c,
		voters: voters,
	}

	s := grpc.NewServer()
	tm.Register(s)
	admin.Register(s, r, cn.voters)

	cs.RegisterCacheServiceServer(s, &cacheServer{c: cn})

	_, port, err := net.SplitHostPort(raftAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse port from %q: %v", raftAddr, err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve raft node at %s: %v", port, err)
		}
	}()

	return cn, nil
}

func newRaftNode(raftID string, raftAddr string, raftDataDir string, bootstrap bool, fsm raft.FSM) (*raft.Raft, *transport.Manager, error) {
	// Create a log store and stable store.
	logStore, err := raftboltdb.NewBoltStore(filepath.Join(raftDataDir, "logs.dat"))
	if err != nil {
		return nil, nil, fmt.Errorf(`boltdb.NewBoltStore(%q): %v`, filepath.Join(raftDataDir, "logs.dat"), err)
	}

	stableStore, err := raftboltdb.NewBoltStore(filepath.Join(raftDataDir, "stable.dat"))
	if err != nil {
		return nil, nil, fmt.Errorf(`boltdb.NewBoltStore(%q): %v`, filepath.Join(raftDataDir, "stable.dat"), err)
	}

	// Create the snapshot store.
	snapshotStore, err := raft.NewFileSnapshotStore(raftDataDir, 2, os.Stderr)
	if err != nil {
		return nil, nil, fmt.Errorf(`raft.NewFileSnapshotStore(%q): %v`, raftDataDir, err)
	}

	// Create a transport layer.
	tm := transport.New(raft.ServerAddress(raftAddr), []grpc.DialOption{grpc.WithInsecure()})

	// Create the configuration for the Raft server.
	raftConfig := raft.DefaultConfig()
	raftConfig.LocalID = raft.ServerID(raftID)

	// Create the Raft server.
	r, err := raft.NewRaft(raftConfig, fsm, logStore, stableStore, snapshotStore, tm.Transport())
	if err != nil {
		return nil, nil, fmt.Errorf("raft.NewRaft: %v", err)
	}

	if bootstrap {
		cfg := raft.Configuration{
			Servers: []raft.Server{
				{
					Suffrage: raft.Voter,
					ID:       raft.ServerID(raftID),
					Address:  raft.ServerAddress(raftAddr),
				},
			},
		}
		log.Printf("Bootstrapping cluster with configuration: %v", cfg)
		f := r.BootstrapCluster(cfg)
		if err := f.Error(); err != nil {
			return nil, nil, fmt.Errorf("error bootstrapping cluster: %v", err)
		}
	}

	return r, tm, nil
}

func checkForAutomaticBootstrap(id string, raftDataDir string) bool {
	_, err := os.Stat(filepath.Join(raftDataDir, "logs.dat"))
	return strings.HasSuffix(id, "-0") && os.IsNotExist(err)
}

func (c *cacheNode) Get(bucket string, key string) (string, bool) {
	value, ok := c.c.Get(bucket, key)
	return value, ok
}

func (c *cacheNode) Put(bucket string, key string, value string) error {
	if c.r.State() != raft.Leader {
		log.Println("Not leader, forwarding request to leader")
		//get leader address
		leader := c.r.Leader()

		conn, err := grpc.Dial(string(leader), grpc.WithInsecure())
		if err != nil {
			return err
		}
		defer conn.Close()
		client := cs.NewCacheServiceClient(conn)

		_, err = client.Put(context.Background(), &cs.PutRequest{
			Bucket: bucket,
			Key:    key,
			Value:  value,
		})
		if err != nil {
			return fmt.Errorf("error forwarding request to leader: %v", err)
		}

		return nil
	}

	cCmd := command{
		Op:     OpSet,
		Bucket: bucket,
		Key:    key,
		Value:  value,
	}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(&cCmd); err != nil {
		return fmt.Errorf("error encoding key, value: %v", err)
	}

	f := c.r.Apply(buf.Bytes(), time.Second)
	if err := f.Error(); err != nil {
		return fmt.Errorf("error applying command to Raft log: %v", err)
	}
	return nil
}

func (c *cacheNode) CreateBucket(bucket string) error {
	if c.r.State() != raft.Leader {
		log.Println("Not leader, forwarding request to leader")
		//get leader address
		leader := c.r.Leader()

		conn, err := grpc.Dial(string(leader), grpc.WithInsecure())
		if err != nil {
			return err
		}
		defer conn.Close()
		client := cs.NewCacheServiceClient(conn)

		_, err = client.CreateBucket(context.Background(), &cs.CreateBucketRequest{
			Bucket: bucket,
		})
		if err != nil {
			return fmt.Errorf("error forwarding request to leader: %v", err)
		}

		return nil
	}

	cCmd := command{
		Op:     OpCreateBucket,
		Bucket: bucket,
	}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(&cCmd); err != nil {
		return fmt.Errorf("error encoding key, value: %v", err)
	}

	f := c.r.Apply(buf.Bytes(), time.Second)
	if err := f.Error(); err != nil {
		return fmt.Errorf("error applying command to Raft log: %v", err)
	}
	return nil
}
