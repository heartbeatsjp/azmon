package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

const (
	// Name is the application name
	Name = "azmon"
	// Usage is the application usage
	Usage = "A tool for Azure Monitor at Microsoft Azure"
	// Version is the application version
	Version = "0.0.0"
)

func buildFetchMetricDataInput(c *cli.Context) FetchMetricDataInput {
	//subCommand := strings.Split(c.Args().First(), " ")[0]
	subCommand := c.Parent().Args().First()

	var metricNames []string
	if subCommand == "check" {
		metricNames = []string{c.String("metric-name")}
	} else if subCommand == "metric" {
		metricNames = c.StringSlice("metric-names")
	}

	return FetchMetricDataInput{
		subscriptionID: c.GlobalString("subscription-id"),
		resourceGroup:  c.GlobalString("resource-group"),
		namespace:      c.GlobalString("namespace"),
		resource:       c.GlobalString("resource"),
		metricNames:    metricNames,
		aggregation:    c.GlobalString("aggregation"),
	}
}

func buildFetchMetricDefinitionsInput(c *cli.Context) FetchMetricDefinitionsInput {
	return FetchMetricDefinitionsInput{
		subscriptionID: c.GlobalString("subscription-id"),
		resourceGroup:  c.GlobalString("resource-group"),
		namespace:      c.GlobalString("namespace"),
		resource:       c.GlobalString("resource"),
	}
}

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
	app.Name = Name
	app.Usage = Usage
	app.Version = Version

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
					Name:  "metric-name, m",
					Usage: "Set the name of the metric",
				},
				cli.Float64Flag{
					Name:  "warning-over, w",
					Usage: "Set the warning threshold",
					Value: 0,
				},
				cli.Float64Flag{
					Name:  "warning-under, W",
					Usage: "Set the warning threshold",
					Value: 0,
				},
				cli.Float64Flag{
					Name:  "critical-over, c",
					Usage: "Set the critical threshold",
					Value: 0,
				},
				cli.Float64Flag{
					Name:  "critical-under, C",
					Usage: "Set the critical threshold",
					Value: 0,
				},
			},
			Action: Check,
		},
		{
			Name:  "metric",
			Usage: "list metric data",
			Flags: []cli.Flag{
				cli.StringSliceFlag{
					Name:  "metric-names, m",
					Usage: "Set the names of the metric",
				},
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
