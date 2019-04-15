package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/urfave/cli"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeError = 10 + iota
)

// Metric prints metric data
func Metric(c *cli.Context) error {
	dataInput := buildFetchMetricDataInput(c)
	defInput := buildFetchMetricDefinitionsInput(c)

	client, err := NewClient(c.GlobalString("subscription-id"))
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("construct api client failed: %s", err.Error()), ExitCodeError)
	}

	output, merr := _metric(client, dataInput, defInput, c.String("prefix"))
	if merr != nil {
		return merr
	}

	fmt.Print(output)

	return nil
}

func _metric(client *Client, dataInput FetchMetricDataInput, defInput FetchMetricDefinitionsInput, prefix string) (string, error) {
	if len(dataInput.metricNames) < 1 {
		definitions, err := FetchMetricDefinitions(context.TODO(), client, defInput)
		if err != nil {
			return "", cli.NewExitError(fmt.Sprintf("fetch metric definitions failed: %s", err.Error()), ExitCodeError)
		}

		for _, d := range *definitions {
			dataInput.metricNames = append(dataInput.metricNames, *d.Name.Value)
		}
	}

	metrics, err := FetchMetricData(context.TODO(), client, dataInput)
	if err != nil {
		return "", cli.NewExitError(fmt.Sprintf("fetch metric data failed: %s", err.Error()), ExitCodeError)
	}

	var output string
	for k, v := range metrics {
		metricKey := strings.Replace(
			strings.Join(
				[]string{
					prefix,
					dataInput.namespace,
					dataInput.resourceGroup,
					dataInput.resource,
					k,
					dataInput.aggregation,
				},
				".",
			),
			"/", "", -1,
		)
		metricKey = strings.Replace(metricKey, " ", "", -1)

		var data float64
		switch dataInput.aggregation {
		case "Total":
			data = *v.Total
		case "Average":
			data = *v.Average
		case "Maximum":
			data = *v.Maximum
		case "Minimum":
			data = *v.Minimum
		}

		output += fmt.Sprintf("%s\t%f\t%d\n", metricKey, data, v.TimeStamp.Unix())
	}

	return output, nil
}
