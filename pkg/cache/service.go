package cache

import (
	"context"

	"example.com/raft-cache/proto/cs"
)

type cacheServer struct {
	cs.UnimplementedCacheServiceServer
	c Cache
}

func (s *cacheServer) Put(ctx context.Context, req *cs.PutRequest) (*cs.PutResponse, error) {
	err := s.c.Put(req.Key, req.Value)
	if err != nil {
		return nil, err
	}

	return &cs.PutResponse{
		Success: true,
	}, nil
}
