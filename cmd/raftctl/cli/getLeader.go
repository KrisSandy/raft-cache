package cli

import (
	"context"
	"flag"
	"fmt"
	"os"

	"example.com/raft-cache/proto/ca"

	"example.com/raft-cache/cmd/raftctl/utils"

	"github.com/mitchellh/cli"
	"google.golang.org/grpc"
)

type getLeaderCommand struct {
	addr string
}

func NewGetLeaderCommand() cli.Command {
	return utils.AdapterCommand(&getLeaderCommand{})
}

func (c *getLeaderCommand) Help() string {
	helpText := `
	Usage: raftctl get-leader --addr=<address>
	`

	return helpText
}

func (c *getLeaderCommand) Name() string {
	return "get-leader"
}

func (*getLeaderCommand) Synopsis() string {
	return "Get the leader node in the cluster"
}

func (c *getLeaderCommand) AppendFlags(f *flag.FlagSet) {
	f.StringVar(&c.addr, "addr", os.Getenv("RAFT_ADDR"), "Address of the node")
}

func (c *getLeaderCommand) Run() error {
	if err := c.validate(); err != nil {
		return err
	}

	conn, err := grpc.Dial(c.addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := ca.NewCacheAdminClient(conn)

	resp, err := client.GetLeader(context.Background(), &ca.LeaderRequest{})
	if err != nil {
		return fmt.Errorf("failed to get leader: %v", err)
	}

	fmt.Printf("leader id: %s, address : %s\n", resp.GetId(), resp.GetAddress())

	return nil
}

func (c *getLeaderCommand) validate() error {
	if c.addr == "" {
		return fmt.Errorf("--addr flag or RAFT_ADDR env variable must be set")
	}

	return nil
}
