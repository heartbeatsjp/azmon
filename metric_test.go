package main

import (
	"strings"
	"testing"
)

func Test_metric(t *testing.T) {
	dataInput := FetchMetricDataInput{
		aggregation:   "Average",
		namespace:     "Microsoft.Compute/virtualMachines",
		resourceGroup: "testrg",
		resource:      "testvm",
	}

	defInput := FetchMetricDefinitionsInput{}

	want := []string{
		"azure.Microsoft.ComputevirtualMachines.testrg.testvm.PercentageCPU.Average	10.000000	1550223420",
		"azure.Microsoft.ComputevirtualMachines.testrg.testvm.NetworkIn.Average	10000.000000	1550223420",
		"azure.Microsoft.ComputevirtualMachines.testrg.testvm.NetworkOut.Average	1000.000000	1550223420",
	}

	stdout, stderr, err := _metric(NewFakeClient(getFakeData()), dataInput, defInput, "azure")
	if err != nil {
		t.Error(err)
	}

	for _, w := range want {
		if !strings.Contains(stdout, w) {
			t.Errorf("%s is not contains output", w)
		}
	}

	if stderr != "" {
		t.Error("stderr is not empty")
	}
}
