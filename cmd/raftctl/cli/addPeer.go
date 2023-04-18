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

type addPeerCommand struct {
	id     string
	addr   string
	leader string
}

func NewAddPeerCommand() cli.Command {
	return utils.AdapterCommand(&addPeerCommand{})
}

func (c *addPeerCommand) Help() string {
	helpText := `
	Usage: raftctl add-peer
	`

	return helpText
}

func (c *addPeerCommand) Name() string {
	return "add-peer"
}

func (*addPeerCommand) Synopsis() string {
	return "Add a new peer to the cluster"
}

func (c *addPeerCommand) AppendFlags(f *flag.FlagSet) {
	f.StringVar(&c.id, "id", os.Getenv("RAFT_ID"), "ID of the new node")
	f.StringVar(&c.addr, "addr", os.Getenv("RAFT_ADDR"), "Address of the new node")
	f.StringVar(&c.leader, "leader", "", "Address of the leader node")
}

func (c *addPeerCommand) Run() error {
	if err := c.validate(); err != nil {
		return err
	}

	ctx := context.Background()

	conn, err := grpc.Dial(c.leader, grpc.WithInsecure())
	if err != nil {
		return err
	}

	client := ca.NewCacheAdminClient(conn)

	lResp, err := client.GetLeader(ctx, &ca.LeaderRequest{})
	if err != nil {
		return fmt.Errorf("failed to get leader: %v", err)
	}
	leaderAddr := lResp.GetAddress()

	cResp, err := client.GetConfiguration(ctx, &ca.ConfigurationRequest{})
	if err != nil {
		return fmt.Errorf("failed to get leader: %v", err)
	}
	isEvenNodes := (len(cResp.GetVoters())+len(cResp.GetNonVoters()))%2 == 0

	if err := conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %v", err)
	}

	conn, err = grpc.Dial(leaderAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}

	client = ca.NewCacheAdminClient(conn)

	if isEvenNodes {
		// make non voters as voters
		fmt.Printf("Based on current cluster configuration, adding new node as voter")
		fmt.Printf("Making existing non voters as voters")
		for _, n := range cResp.GetNonVoters() {
			fmt.Printf("Making %s as voter", n.GetId())
			_, err = client.PromotePeer(ctx, &ca.PromotePeerRequest{
				Id:            n.GetId(),
				Address:       n.GetAddress(),
				PreviousIndex: 0,
			})
			if err != nil {
				return fmt.Errorf("failed to promote voter: %v", err)
			}
		}

		// add new node as voter
		fmt.Printf("Adding new node as voter")
		resp, err := client.AddVoter(ctx, &ca.AddVoterRequest{
			Id:            c.id,
			Address:       c.addr,
			PreviousIndex: 0,
		})
		if err != nil {
			return fmt.Errorf("failed to add voter: %v", err)
		}

		fmt.Printf("New node added as voter successfully: %d", resp.Index)
	} else {
		// add new node as non voter
		fmt.Printf("Based on current cluster configuration, adding new node as non voter")
		resp, err := client.AddNonVoter(ctx, &ca.AddNonVoterRequest{
			Id:            c.id,
			Address:       c.addr,
			PreviousIndex: 0,
		})
		if err != nil {
			return fmt.Errorf("failed to add non voter: %v", err)
		}

		fmt.Printf("New peer added as non voter successfully: %d", resp.Index)
	}

	return nil
}

func (c *addPeerCommand) validate() error {
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
