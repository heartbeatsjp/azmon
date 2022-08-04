package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/urfave/cli"
)

// Return codes that represents service status
const (
	OK       = 0
	WARNING  = 1
	CRITICAL = 2
	UNKNOWN  = 3
)

// Check checks metric data
func Check(c *cli.Context) error {
	client, err := NewClient(c.GlobalString("subscription-id"))
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("construct api client failed: %s", err.Error()), UNKNOWN)
	}

	if c.String("metric-name") == "" {
		return cli.NewExitError("missing metric-name", UNKNOWN)
	}

	if strings.Contains(c.String("metric-name"), ",") {
		return cli.NewExitError("the metric-name option can only specify a single metric name", UNKNOWN)
	}

	input := buildFetchMetricDataInput(c)

	warningOver := c.Float64("warning-over")
	warningUnder := c.Float64("warning-under")

	criticalOver := c.Float64("critical-over")
	criticalUnder := c.Float64("critical-under")

	return _check(client, input, warningOver, warningUnder, criticalOver, criticalUnder)
}

func _check(client *Client, input FetchMetricDataInput, warningOver, warningUnder, criticalOver, criticalUnder float64) error {
	metrics, err := FetchMetricData(context.TODO(), client, input)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("fetch metric data failed: %s", err.Error()), UNKNOWN)
	}

	v, ok := metrics[input.metricNames[0]]
	if !ok {
		return cli.NewExitError("UNKNOWN - No datapoint", UNKNOWN)
	}

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

	if criticalOver != 0 && data > criticalOver {
		return cli.NewExitError(fmt.Sprintf("CRITICAL - %s %s is %f that over than %f", input.resource, input.metricNames[0], data, criticalOver), CRITICAL)
	}

	if criticalUnder != 0 && data < criticalUnder {
		return cli.NewExitError(fmt.Sprintf("CRITICAL - %s %s is %f that under than %f", input.resource, input.metricNames[0], data, criticalUnder), CRITICAL)
	}

	if warningOver != 0 && data > warningOver {
		return cli.NewExitError(fmt.Sprintf("WARNING - %s %s is %f that over than %f", input.resource, input.metricNames[0], data, warningOver), WARNING)
	}

	if warningUnder != 0 && data < warningUnder {
		return cli.NewExitError(fmt.Sprintf("WARNING - %s %s is %f that under than %f", input.resource, input.metricNames[0], data, warningUnder), WARNING)
	}

	return cli.NewExitError(fmt.Sprintf("OK - %s %s is %f", input.resource, input.metricNames[0], data), OK)
}
