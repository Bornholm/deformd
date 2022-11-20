package command

import (
	"fmt"
	"strings"

	"github.com/Bornholm/deformd/internal/command/common"
	"github.com/Bornholm/deformd/internal/server"
	"github.com/pkg/errors"
	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader"
	"github.com/urfave/cli/v2"
	"gitlab.com/wpetit/goweb/logger"
)

func Run() *cli.Command {
	flags := common.Flags()

	return &cli.Command{
		Name:  "run",
		Usage: "Run the deformd server",
		Flags: flags,
		Action: func(ctx *cli.Context) error {
			conf, err := common.LoadConfig(ctx)
			if err != nil {
				return errors.Wrap(err, "Could not load configuration")
			}

			logger.SetFormat(logger.Format(conf.Logger.Format))
			logger.SetLevel(logger.Level(conf.Logger.Level))

			srv := server.New(
				server.WithConfig(conf),
			)

			addrs, srvErrs := srv.Start(ctx.Context)

			url := fmt.Sprintf("http://%s", (<-addrs).String())
			url = strings.Replace(url, "0.0.0.0", "127.0.0.1", 1)

			logger.Info(ctx.Context, "listening", logger.F("url", url))

			if err = <-srvErrs; err != nil {
				return errors.WithStack(err)
			}

			return nil
		},
	}
}
