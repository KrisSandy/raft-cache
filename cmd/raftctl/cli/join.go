package cli

import (
	"context"
	"flag"
	"fmt"
	"os"

	"example.com/raft-cache/cmd/raftctl/utils"
	"example.com/raft-cache/proto/ca"
	"github.com/mitchellh/cli"
	"google.golang.org/grpc"
)

type joinCommand struct {
	id     string
	addr   string
	leader string
}

func NewJoinCommand() cli.Command {
	return utils.AdapterCommand(&joinCommand{})
}

func (c *joinCommand) Help() string {
	helpText := `
	Usage: raftctl join
	`

	return helpText
}

func (c *joinCommand) Name() string {
	return "join"
}

func (*joinCommand) Synopsis() string {
	return "Join an existing cluster"
}

func (c *joinCommand) AppendFlags(f *flag.FlagSet) {
	f.StringVar(&c.id, "id", os.Getenv("RAFT_ID"), "ID of the new node")
	f.StringVar(&c.addr, "addr", os.Getenv("RAFT_ADDR"), "Address of the new node")
	f.StringVar(&c.leader, "leader", "", "Address of the leader node")
}

func (c *joinCommand) Run() error {
	if err := c.validate(); err != nil {
		return err
	}

	if c.addr == c.leader {
		return nil
	}

	for {
		conn, err := grpc.Dial(c.addr, grpc.WithInsecure())
		if err != nil {
			fmt.Printf("failed to connect to addr: %v", err)
			continue
		}
		defer conn.Close()

		client := ca.NewCacheAdminClient(conn)

		resp, err := client.GetState(context.Background(), &ca.StateRequest{})
		if err != nil {
			fmt.Printf("failed to get state: %v", err)
			continue
		}

		if resp.State == "follower" {
			break
		}
	}

	conn, err := grpc.Dial(c.leader, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("failed to connect to leader: %v", err)
	}
	defer conn.Close()

	client := ca.NewCacheAdminClient(conn)

	resp, err := client.Join(context.Background(), &ca.JoinRequest{
		Id:            c.id,
		Address:       c.addr,
		PreviousIndex: 0,
	})
	if err != nil {
		return fmt.Errorf("failed to join: %v", err)
	}

	fmt.Printf("%s added successfully: %d", resp.State, resp.Index)

	return nil
}

func (c *joinCommand) validate() error {
	if c.id == "" {
		return fmt.Errorf("--id flag or RAFT_ID env variable must be set")
	}

	if c.addr == "" {
		return fmt.Errorf("--addr flag or RAFT_ADDR env variable must be set")
	}

	if c.leader == "" {
		return fmt.Errorf("--leader is required")
	}

	return nil
}
