package command

import (
	"github.com/Bornholm/deformd/internal/command/config"
	"github.com/urfave/cli/v2"
)

func Root() []*cli.Command {
	return []*cli.Command{
		Run(),
		config.Root(),
	}
}
