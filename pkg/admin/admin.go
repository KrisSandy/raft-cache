package admin

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/raft"
	"google.golang.org/grpc"

	"example.com/raft-cache/proto/ca"
)

type cacheAdminServer struct {
	ca.UnimplementedCacheAdminServer
	r      *raft.Raft
	voters map[string]struct{}
}

func Register(s *grpc.Server, r *raft.Raft, voters map[string]struct{}) {
	ca.RegisterCacheAdminServer(s, &cacheAdminServer{r: r, voters: voters})
}

func (s *cacheAdminServer) GetLeader(ctx context.Context, req *ca.LeaderRequest) (*ca.LeaderResponse, error) {
	addr, id := s.r.LeaderWithID()

	return &ca.LeaderResponse{
		Address: string(addr),
		Id:      string(id),
	}, nil
}

func (s *cacheAdminServer) GetConfiguration(ctx context.Context, req *ca.ConfigurationRequest) (*ca.ConfigurationResponse, error) {
	configFuture := s.r.GetConfiguration()

	if err := configFuture.Error(); err != nil {
		return nil, err
	}

	config := configFuture.Configuration()

	var voters []*ca.Server
	var nonVoters []*ca.Server
	for _, svr := range config.Servers {
		server := &ca.Server{
			Id:      string(svr.ID),
			Address: string(svr.Address),
		}

		if svr.Suffrage == raft.Voter {
			voters = append(voters, server)
		} else {
			nonVoters = append(nonVoters, server)
		}
	}

	return &ca.ConfigurationResponse{
		Voters:    voters,
		NonVoters: nonVoters,
	}, nil
}

func (s *cacheAdminServer) GetState(ctx context.Context, req *ca.StateRequest) (*ca.StateResponse, error) {
	state := s.r.State()

	var stateStr string
	switch state {
	case raft.Leader:
		stateStr = "leader"
	case raft.Follower:
		stateStr = "follower"
	case raft.Candidate:
		stateStr = "candidate"
	}

	return &ca.StateResponse{
		State: stateStr,
	}, nil
}

func (s *cacheAdminServer) Join(ctx context.Context, req *ca.JoinRequest) (*ca.JoinResponse, error) {
	if _, ok := s.voters[req.Id]; ok {
		f := s.r.AddVoter(raft.ServerID(req.Id), raft.ServerAddress(req.Address), 0, time.Second)

		if err := f.Error(); err != nil {
			return nil, err
		}

		return &ca.JoinResponse{
			Index: f.Index(),
			State: "voter",
		}, nil
	} else {
		f := s.r.AddNonvoter(raft.ServerID(req.Id), raft.ServerAddress(req.Address), 0, time.Second)
		if err := f.Error(); err != nil {
			return nil, err
		}

		return &ca.JoinResponse{
			Index: f.Index(),
			State: "non-voter",
		}, nil
	}
}

func (s *cacheAdminServer) AddVoter(ctx context.Context, req *ca.AddVoterRequest) (*ca.AddVoterResponse, error) {
	res := s.r.AddVoter(raft.ServerID(req.Id), raft.ServerAddress(req.Address), req.PreviousIndex, time.Second)

	if res.Error() != nil {
		return nil, res.Error()
	}

	return &ca.AddVoterResponse{
		Index: res.Index(),
	}, nil
}

func (s *cacheAdminServer) AddNonVoter(ctx context.Context, req *ca.AddNonVoterRequest) (*ca.AddNonVoterResponse, error) {
	res := s.r.AddNonvoter(raft.ServerID(req.Id), raft.ServerAddress(req.Address), req.PreviousIndex, time.Second)

	if res.Error() != nil {
		return nil, res.Error()
	}

	return &ca.AddNonVoterResponse{
		Index: res.Index(),
	}, nil
}

func (s *cacheAdminServer) RemovePeer(ctx context.Context, req *ca.RemovePeerRequest) (*ca.RemovePeerResponse, error) {
	res := s.r.RemoveServer(raft.ServerID(req.Id), req.PreviousIndex, time.Second)

	if res.Error() != nil {
		return nil, res.Error()
	}

	return &ca.RemovePeerResponse{
		Index: res.Index(),
	}, nil
}

func (s *cacheAdminServer) PromotePeer(ctx context.Context, req *ca.PromotePeerRequest) (*ca.PromotePeerResponse, error) {
	_, err := s.RemovePeer(ctx, &ca.RemovePeerRequest{
		Id:            req.Id,
		PreviousIndex: req.PreviousIndex,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to remove server: %w", err)
	}

	avRes, err := s.AddVoter(ctx, &ca.AddVoterRequest{
		Id:            req.Id,
		Address:       req.Address,
		PreviousIndex: req.PreviousIndex,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add voter: %w", err)
	}

	return &ca.PromotePeerResponse{
		Index: avRes.Index,
	}, nil
}
