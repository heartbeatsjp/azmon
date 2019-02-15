package main

import (
	"context"
	"fmt"

	"github.com/urfave/cli"
)

const (
	OK       = 0
	WARNING  = 1
	CRITICAL = 2
	UNKNOWN  = 3
)

func Check(c *cli.Context) error {
	input := buildFetchMetricDataInput(c)

	warningOver := c.Float64("warning-over")
	warningUnder := c.Float64("warning-under")

	criticalOver := c.Float64("critical-over")
	criticalUnder := c.Float64("critical-under")

	v, err := FetchMetricData(context.TODO(), input)
	if err != nil {
		return cli.NewExitError(fmt.Errorf("fetch metric data failed: %s", err.Error()), UNKNOWN)
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
		return cli.NewExitError(fmt.Sprintf("CRITICAL - %s %s is %f that over than %f", input.resource, input.metricName, data, criticalOver), CRITICAL)
	}

	if criticalUnder != 0 && data < criticalUnder {
		return cli.NewExitError(fmt.Sprintf("CRITICAL - %s %s is %f that under than %f", input.resource, input.metricName, data, criticalUnder), CRITICAL)
	}

	if warningOver != 0 && data > warningOver {
		return cli.NewExitError(fmt.Sprintf("WARNING - %s %s is %f that over than %f", input.resource, input.metricName, data, warningOver), WARNING)
	}

	if warningUnder != 0 && data < warningUnder {
		return cli.NewExitError(fmt.Sprintf("WARNING - %s %s is %f that under than %f", input.resource, input.metricName, data, warningUnder), WARNING)
	}

	fmt.Printf("OK - %s %s is %f", input.resource, input.metricName, data)
	return nil
}
