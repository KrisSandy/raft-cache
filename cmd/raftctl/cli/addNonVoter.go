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

type addNonVoterCommand struct {
	id     string
	addr   string
	leader string
}

func NewAddNonVoterCommand() cli.Command {
	return utils.AdapterCommand(&addNonVoterCommand{})
}

func (c *addNonVoterCommand) Help() string {
	helpText := `
	Usage: raftctl add-non-voter
	`

	return helpText
}

func (c *addNonVoterCommand) Name() string {
	return "add-non-voter"
}

func (*addNonVoterCommand) Synopsis() string {
	return "Add a new non-voter to the cluster"
}

func (c *addNonVoterCommand) AppendFlags(f *flag.FlagSet) {
	f.StringVar(&c.id, "id", os.Getenv("RAFT_ID"), "ID of the new node")
	f.StringVar(&c.addr, "addr", os.Getenv("RAFT_ADDR"), "Address of the new node")
	f.StringVar(&c.leader, "leader", "", "Address of the leader node")

}

func (c *addNonVoterCommand) Run() error {
	if err := c.validate(); err != nil {
		return err
	}

	conn, err := grpc.Dial(c.leader, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := ca.NewCacheAdminClient(conn)

	resp, err := client.AddNonVoter(context.Background(), &ca.AddNonVoterRequest{
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

func (c *addNonVoterCommand) validate() error {
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
