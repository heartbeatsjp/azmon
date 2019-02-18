package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/urfave/cli"
)

const (
	OK       = 0
	WARNING  = 1
	CRITICAL = 2
	UNKNOWN  = 3
)

func Check(c *cli.Context) error {
	if c.String("metric-name") == "" {
		return cli.NewExitError("missing metric-name", UNKNOWN)
	}

	if strings.Contains(c.String("metric-name"), ",") {
		//TODO: error message
		return cli.NewExitError("TODO", UNKNOWN)
	}

	input := buildFetchMetricDataInput(c)

	warningOver := c.Float64("warning-over")
	warningUnder := c.Float64("warning-under")

	criticalOver := c.Float64("critical-over")
	criticalUnder := c.Float64("critical-under")

	client, err := NewClient(input.subscriptionID)
	if err != nil {
		return cli.NewExitError("", UNKNOWN)
	}
	metrics, err := FetchMetricData(context.TODO(), client, input)
	if err != nil {
		return cli.NewExitError(fmt.Errorf("fetch metric data failed: %s", err.Error()), UNKNOWN)
	}

	v := metrics[input.metricNames[0]]

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
