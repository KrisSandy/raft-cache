package cli

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/mitchellh/cli"
	"google.golang.org/grpc"

	"example.com/raft-cache/proto/ca"

	"example.com/raft-cache/cmd/raftctl/utils"
)

type listPeersCommand struct {
	leader string
}

func NewListCommand() cli.Command {
	return utils.AdapterCommand(&listPeersCommand{})
}

func (c *listPeersCommand) Help() string {
	helpText := `
	Usage: raftctl list-peers --leader=<address>
	`

	return helpText
}

func (c *listPeersCommand) Name() string {
	return "list-peers"
}

func (*listPeersCommand) Synopsis() string {
	return "List all peers in the cluster"
}

func (c *listPeersCommand) AppendFlags(f *flag.FlagSet) {
	f.StringVar(&c.leader, "leader", os.Getenv("RAFT_ADDR"), "Address of the peer node")
}

func (c *listPeersCommand) Run() error {
	if err := c.validate(); err != nil {
		return err
	}

	ctx := context.Background()
	conn, err := grpc.Dial(c.leader, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := ca.NewCacheAdminClient(conn)

	resp, err := client.GetConfiguration(ctx, &ca.ConfigurationRequest{})
	if err != nil {
		return fmt.Errorf("failed to get leader: %v", err)
	}

	fmt.Println("VOTERS:")
	printServer(ctx, client, resp.GetVoters())

	fmt.Println("\nNON-VOTERS:")
	printServer(ctx, client, resp.GetNonVoters())

	return nil
}

func (c *listPeersCommand) validate() error {
	if c.leader == "" {
		return fmt.Errorf("--leader flag or RAFT_ADDR env variable must be set")
	}

	return nil
}

func printServer(ctx context.Context, client ca.CacheAdminClient, servers []*ca.Server) error {
	for _, server := range servers {
		state, err := getState(ctx, server)
		if err != nil {
			return err
		}
		fmt.Printf("voter id: %s, address: %s, state: %s\n", server.GetId(), server.GetAddress(), *state)
	}

	return nil
}

func getState(ctx context.Context, server *ca.Server) (*string, error) {
	conn, err := grpc.Dial(server.Address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := ca.NewCacheAdminClient(conn)

	resp, err := client.GetState(ctx, &ca.StateRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get state: %v", err)
	}

	state := resp.GetState()

	return &state, nil
}
