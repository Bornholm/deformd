package common

import "github.com/urfave/cli/v2"

func Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:      "config",
			Aliases:   []string{"c"},
			EnvVars:   []string{"DEFORMD_CONFIG"},
			Value:     "",
			TakesFile: true,
		},
	}
}
