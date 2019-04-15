package main

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-09-01/insights"
	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-09-01/insights/insightsapi"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/Azure/go-autorest/autorest/to"
)

type fakeMetricsClient struct {
	insightsapi.MetricsClientAPI
}

type fakeMetricDefinitionsClient struct {
	insightsapi.MetricDefinitionsClientAPI
}

func (c *fakeMetricsClient) List(ctx context.Context, resourceURI string, timespan string, interval *string, metricnames string, aggregation string, top *int32, orderby string, filter string, resultType insights.ResultType, metricnamespace string) (insights.Response, error) {
	t := time.Unix(1550223420, 0)

	return insights.Response{
		Value: &[]insights.Metric{
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
		},
	}, nil
}

func (c *fakeMetricDefinitionsClient) List(ctx context.Context, resourceURI string, metricnamespace string) (insights.MetricDefinitionCollection, error) {
	return insights.MetricDefinitionCollection{
		Value: &[]insights.MetricDefinition{
			{Name: &insights.LocalizableString{Value: to.StringPtr("Percentage CPU")}, IsDimensionRequired: to.BoolPtr(false)},
			{Name: &insights.LocalizableString{Value: to.StringPtr("Network In")}, IsDimensionRequired: to.BoolPtr(false)},
			{Name: &insights.LocalizableString{Value: to.StringPtr("Network Out")}, IsDimensionRequired: to.BoolPtr(false)},
		},
	}, nil
}

func NewFakeClient() *Client {
	return &Client{
		MetricsClient:           &fakeMetricsClient{},
		MetricDefinitionsClient: &fakeMetricDefinitionsClient{},
	}
}
