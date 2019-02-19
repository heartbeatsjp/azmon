package main

import (
	"testing"

	"github.com/urfave/cli"
)

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
				client: NewFakeClient(),
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
				client: NewFakeClient(),
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
				client: NewFakeClient(),
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
				client: NewFakeClient(),
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
				client: NewFakeClient(),
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
				client: NewFakeClient(),
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
