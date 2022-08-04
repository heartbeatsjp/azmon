package main

import (
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-09-01/insights"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/urfave/cli"
)

func getFakeData() *[]insights.Metric {
	t := time.Unix(1550223420, 0)
	return &[]insights.Metric{
		{
			Name: &insights.LocalizableString{Value: to.StringPtr("Percentage CPU")},
			Timeseries: &[]insights.TimeSeriesElement{
				{
					Data: &[]insights.MetricValue{
						{TimeStamp: &date.Time{Time: t}, Average: to.Float64Ptr(10.000000)},
						{TimeStamp: &date.Time{Time: t.Add(-1 * time.Minute)}, Average: to.Float64Ptr(11.000000)},
						{TimeStamp: &date.Time{Time: t.Add(-2 * time.Minute)}, Average: to.Float64Ptr(12.000000)},
						{TimeStamp: &date.Time{Time: t.Add(-3 * time.Minute)}, Average: to.Float64Ptr(13.000000)},
						{TimeStamp: &date.Time{Time: t.Add(-4 * time.Minute)}, Average: to.Float64Ptr(14.000000)},
					},
				},
			},
		},
		{
			Name: &insights.LocalizableString{Value: to.StringPtr("Network In")},
			Timeseries: &[]insights.TimeSeriesElement{
				{
					Data: &[]insights.MetricValue{
						{TimeStamp: &date.Time{Time: t}, Average: to.Float64Ptr(10000.000000)},
						{TimeStamp: &date.Time{Time: t.Add(-1 * time.Minute)}, Average: to.Float64Ptr(11000.000000)},
						{TimeStamp: &date.Time{Time: t.Add(-2 * time.Minute)}, Average: to.Float64Ptr(12000.000000)},
						{TimeStamp: &date.Time{Time: t.Add(-3 * time.Minute)}, Average: to.Float64Ptr(13000.000000)},
						{TimeStamp: &date.Time{Time: t.Add(-4 * time.Minute)}, Average: to.Float64Ptr(14000.000000)},
					},
				},
			},
		},
		{
			Name: &insights.LocalizableString{Value: to.StringPtr("Network Out")},
			Timeseries: &[]insights.TimeSeriesElement{
				{
					Data: &[]insights.MetricValue{
						{TimeStamp: &date.Time{Time: t}, Average: to.Float64Ptr(1000.000000)},
						{TimeStamp: &date.Time{Time: t.Add(-1 * time.Minute)}, Average: to.Float64Ptr(1100.000000)},
						{TimeStamp: &date.Time{Time: t.Add(-2 * time.Minute)}, Average: to.Float64Ptr(1200.000000)},
						{TimeStamp: &date.Time{Time: t.Add(-3 * time.Minute)}, Average: to.Float64Ptr(1300.000000)},
						{TimeStamp: &date.Time{Time: t.Add(-4 * time.Minute)}, Average: to.Float64Ptr(1400.000000)},
					},
				},
			},
		},
	}
}

func Test_check(t *testing.T) {
	type args struct {
		client        *Client
		input         FetchMetricDataInput
		warningOver   float64
		warningUnder  float64
		criticalOver  float64
		criticalUnder float64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Percentage CPU: OK",
			args: args{
				client: NewFakeClient(getFakeData()),
				input: FetchMetricDataInput{
					metricNames: []string{"Percentage CPU"},
					aggregation: "Average",
				},
				warningOver:  10,
				criticalOver: 11,
			},
			want: OK,
		},
		{
			name: "Percentage CPU: WARNING",
			args: args{
				client: NewFakeClient(getFakeData()),
				input: FetchMetricDataInput{
					metricNames: []string{"Percentage CPU"},
					aggregation: "Average",
				},
				warningOver:  9,
				criticalOver: 11,
			},
			want: WARNING,
		},
		{
			name: "Percentage CPU: CRITICAL",
			args: args{
				client: NewFakeClient(getFakeData()),
				input: FetchMetricDataInput{
					metricNames: []string{"Percentage CPU"},
					aggregation: "Average",
				},
				warningOver:  9,
				criticalOver: 9.9,
			},
			want: CRITICAL,
		},
		{
			name: "Network In: OK",
			args: args{
				client: NewFakeClient(getFakeData()),
				input: FetchMetricDataInput{
					metricNames: []string{"Network In"},
					aggregation: "Average",
				},
				warningUnder:  10000,
				criticalUnder: 9000,
			},
			want: OK,
		},
		{
			name: "Network In: WARNING",
			args: args{
				client: NewFakeClient(getFakeData()),
				input: FetchMetricDataInput{
					metricNames: []string{"Network In"},
					aggregation: "Average",
				},
				warningUnder:  11000,
				criticalUnder: 10000,
			},
			want: WARNING,
		},
		{
			name: "Network In: CRITICAL",
			args: args{
				client: NewFakeClient(getFakeData()),
				input: FetchMetricDataInput{
					metricNames: []string{"Network In"},
					aggregation: "Average",
				},
				warningUnder:  12000,
				criticalUnder: 11000,
			},
			want: CRITICAL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := _check(tt.args.client, tt.args.input, tt.args.warningOver, tt.args.warningUnder, tt.args.criticalOver, tt.args.criticalUnder)
			if err == nil {
				t.Error("error")
			}

			if eerr, ok := err.(*cli.ExitError); ok {
				if eerr.ExitCode() != tt.want {
					t.Errorf("ExitCode want %d, got %d\n", tt.want, eerr.ExitCode())
				}
			} else {
				t.Errorf("Invalid error type")
			}
		})
	}
}

func Test_check_no_datapoint(t *testing.T) {
	t.Run("no datapoint", func(t *testing.T) {
		err := _check(
			NewFakeClient(&[]insights.Metric{}), // empty metric data
			FetchMetricDataInput{
				metricNames: []string{"Network In"},
				aggregation: "Average",
			}, 0, 0, 0, 0)
		if err == nil {
			t.Error("error")
		}
		if eerr, ok := err.(*cli.ExitError); ok {
			if eerr.ExitCode() != UNKNOWN {
				t.Errorf("ExitCode want %d, got %d\n", UNKNOWN, eerr.ExitCode())
			}
			if eerr.Error() != "UNKNOWN - No datapoint" {
				t.Errorf("ExitCode want %s, got %s\n", "no datapoint", eerr.Error())
			}
		} else {
			t.Errorf("Invalid error type")
		}
	})
}
