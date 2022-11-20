package common

import (
	"github.com/Bornholm/deformd/internal/config"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func LoadConfig(ctx *cli.Context) (*config.Config, error) {
	configFile := ctx.String("config")

	var (
		conf *config.Config
		err  error
	)

	if configFile != "" {
		conf, err = config.NewFromFile(configFile)
		if err != nil {
			return nil, errors.Wrapf(err, "Could not load config file '%s'", configFile)
		}
	} else {
		conf = config.NewDefault()
	}

	return conf, nil
}
