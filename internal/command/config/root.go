package config

import "github.com/urfave/cli/v2"

func Root() *cli.Command {
	return &cli.Command{
		Name:  "config",
		Usage: "Config related commands",
		Subcommands: []*cli.Command{
			Dump(),
		},
	}
}
