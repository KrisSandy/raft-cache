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

type addVoterCommand struct {
	id     string
	addr   string
	leader string
}

func NewAddVoterCommand() cli.Command {
	return utils.AdapterCommand(&addVoterCommand{})
}

func (c *addVoterCommand) Help() string {
	helpText := `
	Usage: raftctl add-voter
	`

	return helpText
}

func (c *addVoterCommand) Name() string {
	return "add-voter"
}

func (*addVoterCommand) Synopsis() string {
	return "Add a new voter to the cluster"
}

func (c *addVoterCommand) AppendFlags(f *flag.FlagSet) {
	f.StringVar(&c.id, "id", os.Getenv("RAFT_ID"), "ID of the new node")
	f.StringVar(&c.addr, "addr", os.Getenv("RAFT_ADDR"), "Address of the new node")
	f.StringVar(&c.leader, "leader", "", "Address of the leader node")
}

func (c *addVoterCommand) Run() error {
	if err := c.validate(); err != nil {
		return err
	}

	conn, err := grpc.Dial(c.leader, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := ca.NewCacheAdminClient(conn)

	resp, err := client.AddVoter(context.Background(), &ca.AddVoterRequest{
		Id:            c.id,
		Address:       c.addr,
		PreviousIndex: 0,
	})
	if err != nil {
		return fmt.Errorf("failed to add voter: %v", err)
	}

	fmt.Printf("voter added successfully: %d", resp.Index)

	return nil
}

func (c *addVoterCommand) validate() error {
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
