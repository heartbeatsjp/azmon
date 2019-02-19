package main

import (
	"strings"
	"testing"
)

func Test_metric(t *testing.T) {
	dataInput := FetchMetricDataInput{
		aggregation: "Average",
		namespace:   "Microsoft.Compute/virtualMachines",
		resource:    "testvm",
	}

	defInput := FetchMetricDefinitionsInput{}

	want := []string{
		"azure.Microsoft.ComputevirtualMachines.testvm.PercentageCPU.Average	10.000000	1550223420",
		"azure.Microsoft.ComputevirtualMachines.testvm.NetworkIn.Average	10000.000000	1550223420",
		"azure.Microsoft.ComputevirtualMachines.testvm.NetworkOut.Average	1000.000000	1550223420",
	}

	got, err := _metric(NewFakeClient(), dataInput, defInput, "azure")
	if err != nil {
		t.Error(err)
	}

	for _, w := range want {
		if !strings.Contains(got, w) {
			t.Errorf("%s is not contains output", w)
		}
	}
}
