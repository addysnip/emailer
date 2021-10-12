package main

import (
	"fmt"
	"os"

	"addysnip.dev/emailer/cmd/consumer"
	"addysnip.dev/emailer/cmd/migrate"
	"addysnip.dev/emailer/pkg/logger"
	"addysnip.dev/emailer/pkg/version"
	"github.com/common-nighthawk/go-figure"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:                 "addysnip",
		Usage:                "Addysnip.io Eamiler Service",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			consumer.Command(),
			migrate.Command(),
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Print the version",
				Action: func(c *cli.Context) error {
					logger.Category("main").Info("Version: %s", version.FriendlyVersion())
					return nil
				},
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Usage:   "Set logging to verbose",
				Aliases: []string{"v"},
			},
		},
		Before: func(c *cli.Context) error {
			if c.Bool("verbose") {
				logger.Category("cmd").Info("Setting logging to debug")
				logger.SetLogLevel(logger.DEBUG)
			}

			intro := figure.NewFigure("addysnip", "", false).Slicify()
			for i := 0; i < len(intro); i++ {
				fmt.Printf("%s\n", intro[i])
			}
			fmt.Printf("Addysnip Emailer Service %s\n", version.FriendlyVersion())
			logger.Category("main").Info("Checking for .env, if exists, will load")
			if _, err := os.Stat(".env"); err == nil {
				logger.Category("main").Debug("Loading .env")
				err := godotenv.Load()
				if err != nil {
					logger.Category("main").Error("Error loading .env file: %s", err.Error())
					return err
				}
			}

			return nil
		},
	}

	app.Run(os.Args)
}
