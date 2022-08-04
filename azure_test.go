package main

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-09-01/insights"
	"github.com/Azure/go-autorest/autorest/to"
)

type fakeMetricsClient struct {
	fakeData *[]insights.Metric
}

type fakeMetricDefinitionsClient struct {
}

func (c *fakeMetricsClient) List(ctx context.Context, resourceURI string, timespan string, interval *string, metricnames string, aggregation string, top *int32, orderby string, filter string, resultType insights.ResultType, metricnamespace string) (insights.Response, error) {
	return insights.Response{
		Value: c.fakeData,
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

func NewFakeClient(fakeData *[]insights.Metric) *Client {
	return &Client{
		MetricsClient: &fakeMetricsClient{
			fakeData: fakeData,
		},
		MetricDefinitionsClient: &fakeMetricDefinitionsClient{},
	}
}
