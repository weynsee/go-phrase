package cli

import (
	mcli "github.com/mitchellh/cli"
)

type cli struct {
	cli *mcli.CLI
}

// NewCLI returns a CLI instance.
func NewCLI(version string, args []string) *cli {
	c := mcli.NewCLI("go-phrase", version)
	c.Args = args
	c.Commands = commands

	return &cli{c}
}

// Run runs the actual program based on the arguments given.
func (c *cli) Run() (int, error) {
	return c.cli.Run()
}
