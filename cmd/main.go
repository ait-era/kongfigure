package main

import (
	"kongfigure/internal"
	"log"
	"os"
)
import "github.com/urfave/cli"

func main() {
	var settings kongfigure.AppSettings

	app := cli.NewApp()
	app.Name = "Kongfigure"
	app.Usage = "Tool to configure Kong services, routes and plugins."
	app.Version = "0.1.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "kong-configs",
			Usage:       "Path to the folder containing Kong JSON's configuration",
			EnvVar:      "KGF_KONG_PATH",
			Destination: &settings.KongConfPath,
		},
		cli.StringFlag{
			Name:        "kong-url",
			Usage:       "Kong Admin URL",
			EnvVar:      "KGF_KONG_URL",
			Destination: &settings.KongUrl,
		},
		cli.BoolFlag{
			Name:        "dry-run",
			Usage:       "Doesn't actually call Kong and cannot verify if IDs in relation actually exists.",
			EnvVar:      "KGF_DRY_RUN",
			Hidden:      false,
			Destination: &settings.DryRun,
		},
	}

	app.Action = func(c *cli.Context) error {
		return kongfigure.Run(c, settings)
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}

