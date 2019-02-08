package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "subscription-id",
			Usage: "Set the subscription id",
		},
		cli.StringFlag{
			Name:  "resource-group",
			Usage: "Set the resource group name",
		},
		cli.StringFlag{
			Name:  "namespace",
			Usage: "Set the metric namespace",
		},
		cli.StringFlag{
			Name:  "resource",
			Usage: "Set the target resource name",
		},
		cli.StringFlag{
			Name:  "metric-name",
			Usage: "Set the name of the metric",
		},
		cli.StringFlag{
			Name:  "aggregation",
			Usage: "Set the aggregation type",
		},
		cli.StringFlag{
			Name:  "auth-file",
			Usage: "Set the azure auth file path",
			Value: "/etc/nagios/azure.auth",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "check",
			Usage: "check metric(as Nagios plugin)",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "warning,w",
					Usage: "Set the warning threshold",
				},
				cli.StringFlag{
					Name:  "critical,c",
					Usage: "Set the critical threshold",
				},
			},
			Action: func(c *cli.Context) error { return nil },
		},
		{
			Name:   "metric",
			Usage:  "list metric data",
			Action: func(c *cli.Context) error { return nil },
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
