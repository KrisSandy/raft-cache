package cli

import (
	"context"
	"flag"
	"fmt"

	"example.com/raft-cache/proto/ca"

	"example.com/raft-cache/cmd/raftctl/utils"

	"github.com/mitchellh/cli"
	"google.golang.org/grpc"
)

type removePeerCommand struct {
	id     string
	leader string
}

func NewRemovePeerCommand() cli.Command {
	return utils.AdapterCommand(&removePeerCommand{})
}

func (c *removePeerCommand) Help() string {
	helpText := `
	Usage: raftctl remove-peer
	`

	return helpText
}

func (c *removePeerCommand) Name() string {
	return "remove-peer"
}

func (*removePeerCommand) Synopsis() string {
	return "Remove a peer from the cluster"
}

func (c *removePeerCommand) AppendFlags(f *flag.FlagSet) {
	f.StringVar(&c.id, "id", "", "ID of the node")
	f.StringVar(&c.leader, "leader", "", "Address of the leader node")
}

func (c *removePeerCommand) Run() error {
	if err := c.validate(); err != nil {
		return err
	}

	ctx := context.Background()

	conn, err := grpc.Dial(c.leader, grpc.WithInsecure())
	if err != nil {
		return err
	}

	client := ca.NewCacheAdminClient(conn)

	resp, err := client.RemovePeer(ctx, &ca.RemovePeerRequest{
		Id: c.id,
	})
	if err != nil {
		return err
	}

	fmt.Printf("Removed peer %d from the cluster", resp.Index)

	return nil
}

func (c *removePeerCommand) validate() error {
	if c.id == "" {
		return fmt.Errorf("--id is required")
	}

	if c.leader == "" {
		return fmt.Errorf("--leader is required")
	}

	return nil
}
