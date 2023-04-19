package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	RaftID        string
	RaftAddr      string
	RaftDataDir   string
	RaftBootstrap bool
	RaftVoters    map[string]struct{}
	ServicePort   string
}

func NewConfig() (*Config, error) {
	raftId := flag.String("raft-id", os.Getenv("RAFT_ID"), "id of the raft node")
	raftAddr := flag.String("raft-addr", os.Getenv("RAFT_ADDR"), "address of the raft node")
	raftDataDir := flag.String("raft-data-dir", os.Getenv("RAFT_DATA_DIR"), "raft data directory")
	raftBootstrap := flag.Bool("raft-bootstrap", false, "bootstrap the raft cluster")
	raftVoters := flag.String("raft-voters", os.Getenv("RAFT_VOTERS"), "comma separated list of raft peers")
	svrPort := flag.String("svr-port", "8000", "address of the server")

	flag.Parse()

	if *raftId == "" {
		return nil, fmt.Errorf("RAFT_ID is required")
	}

	if *raftAddr == "" {
		return nil, fmt.Errorf("RAFT_ADDR is required")
	}

	if *raftDataDir == "" {
		return nil, fmt.Errorf("RAFT_DATA_DIR is required")
	}

	voters := make(map[string]struct{})
	if *raftVoters != "" {
		vs := strings.Split(*raftVoters, ",")
		for _, v := range vs {
			voters[v] = struct{}{}
		}
	}

	return &Config{
		RaftID:        *raftId,
		RaftAddr:      *raftAddr,
		RaftDataDir:   *raftDataDir,
		RaftBootstrap: *raftBootstrap,
		RaftVoters:    voters,
		ServicePort:   *svrPort,
	}, nil
}
