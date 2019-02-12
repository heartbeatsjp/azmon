package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

func validationGlobalFlags(c *cli.Context) error {
	if v := c.GlobalString("subscription-id"); v == "" {
		return errors.New("missing subscription-id")
	}

	if v := c.GlobalString("resource-group"); v == "" {
		return errors.New("missing resource-group")
	}

	if v := c.GlobalString("namespace"); v == "" {
		return errors.New("missing namespace")
	}

	if v := c.GlobalString("resource"); v == "" {
		return errors.New("missing resource")
	}

	if v := c.GlobalString("metric-name"); v == "" {
		return errors.New("missing metric-name")
	}

	if v := c.GlobalString("aggregation"); v == "" {
		return errors.New("missing aggregation")
	}

	if v := c.GlobalString("aggregation"); v != "Total" && v != "Average" && v != "Maximum" && v != "Minimum" {
		return errors.New("invalid aggregation: choose from \"Total\", \"Average\", \"Maximum\", \"Minimum\" (\"Count\" is not supported)")
	}

	return nil
}

func setAzureAuthLocation(c *cli.Context) error {
	return os.Setenv("AZURE_AUTH_LOCATION", c.GlobalString("auth-file"))
}

func appBefore(c *cli.Context) error {
	if err := validationGlobalFlags(c); err != nil {
		return fmt.Errorf("validation global flags failed: %s", err.Error())
	}
	if err := setAzureAuthLocation(c); err != nil {
		return fmt.Errorf("set AZURE_AUTH_LOCATION failed: %s", err.Error())
	}
	return nil
}

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "subscription-id, s",
			Usage: "Set the subscription id",
		},
		cli.StringFlag{
			Name:  "resource-group, g",
			Usage: "Set the resource group name",
		},
		cli.StringFlag{
			Name:  "namespace, n",
			Usage: "Set the metric namespace",
		},
		cli.StringFlag{
			Name:  "resource, r",
			Usage: "Set the target resource name",
		},
		cli.StringFlag{
			Name:  "metric-name, m",
			Usage: "Set the name of the metric",
		},
		cli.StringFlag{
			Name:  "aggregation, a",
			Usage: "Set the aggregation type. Choose from \"Total\", \"Average\", \"Maximum\", \"Minimum\" (\"Count\" is not supported)",
		},
		cli.StringFlag{
			Name:  "auth-file",
			Usage: "Set the azure auth file path",
			Value: "/etc/nagios/azure.auth",
		},
	}

	app.Before = appBefore

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
			Name:  "metric",
			Usage: "list metric data",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "prefix, p",
					Usage: "Set the metric key prefix",
					Value: "azure",
				},
			},
			Action: Metric,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
