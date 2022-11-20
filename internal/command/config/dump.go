package config

import (
	"os"

	"github.com/Bornholm/deformd/internal/command/common"
	"github.com/Bornholm/deformd/internal/config"
	"github.com/pkg/errors"
	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader"
	"github.com/urfave/cli/v2"
	"gitlab.com/wpetit/goweb/logger"
)

func Dump() *cli.Command {
	flags := common.Flags()

	return &cli.Command{
		Name:  "dump",
		Usage: "Dump the current configuration",
		Flags: flags,
		Action: func(ctx *cli.Context) error {
			conf, err := common.LoadConfig(ctx)
			if err != nil {
				return errors.Wrap(err, "Could not load configuration")
			}

			logger.SetFormat(logger.Format(conf.Logger.Format))
			logger.SetLevel(logger.Level(conf.Logger.Level))

			if err := config.Dump(conf, os.Stdout); err != nil {
				return errors.Wrap(err, "Could not dump configuration")
			}

			return nil
		},
	}
}
