package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/urfave/cli"
)

func Metric(c *cli.Context) error {
	input := buildFetchMetricDataInput(c)

	client, err := NewClient(input.subscriptionID)
	if err != nil {
		return cli.NewExitError("", UNKNOWN)
	}

	if len(input.metricNames) < 1 {
		i := buildFetchMetricDefinitionsInput(c)
		definitions, err := FetchMetricDefinitions(context.TODO(), client, i)
		if err != nil {
			return cli.NewExitError("", 1)
		}

		for _, d := range *definitions {
			input.metricNames = append(input.metricNames, *d.Name.Value)
		}
	}

	metrics, err := FetchMetricData(context.TODO(), client, input)
	if err != nil {
		return fmt.Errorf("fetch metric data failed: %s", err.Error())
	}

	for k, v := range metrics {
		prefix := c.String("prefix")
		metricKey := strings.Replace(
			strings.Join(
				[]string{
					prefix,
					input.namespace,
					input.resource,
					k,
					input.aggregation,
				},
				".",
			),
			"/", "", -1,
		)
		metricKey = strings.Replace(metricKey, " ", "", -1)

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

		fmt.Printf("%s\t%f\t%d\n", metricKey, data, v.TimeStamp.Unix())
	}

	return nil
}
