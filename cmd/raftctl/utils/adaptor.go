package utils

import (
	"flag"
	"fmt"
)

type Command interface {
	Name() string
	Synopsis() string
	Help() string
	AppendFlags(*flag.FlagSet)
	Run() error
}

type Adapter struct {
	cmd Command

	flags *flag.FlagSet
}

func AdapterCommand(cmd Command) *Adapter {
	a := &Adapter{
		cmd: cmd,
	}

	f := flag.NewFlagSet(cmd.Name(), flag.ContinueOnError)
	a.cmd.AppendFlags(f)
	a.flags = f

	return a
}

func (a *Adapter) Help() string {
	return a.cmd.Help()
}

func (a *Adapter) Synopsis() string {
	return a.cmd.Synopsis()
}

func (a *Adapter) Run(args []string) int {
	if err := a.flags.Parse(args); err != nil {
		return 1
	}

	if err := a.cmd.Run(); err != nil {
		fmt.Printf("Error: %s\n", err)
		return 1
	}

	return 0
}
