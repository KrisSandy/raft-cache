package main

import (
	"log"
	"os"

	raftcli "example.com/raft-cache/cmd/raftctl/cli"

	"github.com/mitchellh/cli"
)

// create command line tool to interact with the cache
func main() {
	c := cli.NewCLI("raftctl", "1.0.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"list-peers": func() (cli.Command, error) {
			return raftcli.NewListCommand(), nil
		},
		"get-leader": func() (cli.Command, error) {
			return raftcli.NewGetLeaderCommand(), nil
		},
		"add-voter": func() (cli.Command, error) {
			return raftcli.NewAddVoterCommand(), nil
		},
		"add-non-voter": func() (cli.Command, error) {
			return raftcli.NewAddNonVoterCommand(), nil
		},
		"add-peer": func() (cli.Command, error) {
			return raftcli.NewAddPeerCommand(), nil
		},
		"remove-peer": func() (cli.Command, error) {
			return raftcli.NewRemovePeerCommand(), nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
