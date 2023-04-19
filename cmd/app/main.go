package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"example.com/raft-cache/pkg/config"

	"example.com/raft-cache/pkg/cache"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error creating new config: %v", err)
	}

	log.Printf("Starting cache node %s at %s", config.RaftID, config.RaftAddr)
	c, err := cache.NewCacheNode(config.RaftID, config.RaftAddr, config.RaftDataDir, config.RaftBootstrap, config.RaftVoters)
	if err != nil {
		log.Fatalf("Error creating new raft node: %v", err)
	}

	server := newServer(c)
	log.Printf("Starting server at %s", config.ServicePort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.ServicePort), server.handler))
}

func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
