package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli"
)

func Metric(c *cli.Context) error {
	input := buildFetchMetricDataInput(c)

	if err := os.Setenv("AZURE_AUTH_LOCATION", c.GlobalString("auth-file")); err != nil {
		fmt.Println("set environment variable (AZURE_AUTH_LOCATION) failed")
	}

	v, err := FetchMetricData(
		context.TODO(),
		input.subscriptionID,
		input.resourceGroup,
		input.namespace,
		input.resource,
		input.metricName,
		input.aggregation,
	)
	if err != nil {
		return fmt.Errorf("fetch metric data failed: %s", err.Error())
	}

	prefix := c.String("prefix")
	key := strings.Join(
		[]string{
			prefix,
			strings.Replace(input.namespace, "/", ".", -1),
			input.metricName,
			input.resource,
			input.metricName,
			input.aggregation,
		},
		".",
	)
	key = strings.Replace(key, " ", "", -1)

	var data float64
	switch input.aggregation {
	case "Total":
		data = *v.Total
	case "Average":
		data = *v.Average
	case "Maximum":
		data = *v.Maximum
	case "Minimum":
		data = *v.Minimum
	}

	fmt.Printf("%s\t%f\t%d\n", key, data, v.TimeStamp.Unix())
	return nil
}
