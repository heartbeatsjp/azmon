package main

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-11-01-preview/insights"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
)

type metricsListInput struct {
	subscriptionID  string
	resourceURI     string
	timespan        string
	interval        *string
	metricnames     string
	aggregation     string
	top             *int32
	orderby         string
	filter          string
	resultType      insights.ResultType
	metricnamespace string
}

func metricsList(ctx context.Context, params *metricsListInput) (insights.Response, error) {
	client := insights.NewMetricsClient(params.subscriptionID)
	a, err := auth.NewAuthorizerFromFile(azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		return insights.Response{}, err
	}
	client.Authorizer = a
	return client.List(
		ctx,
		params.resourceURI,
		params.timespan,
		params.interval,
		params.metricnames,
		params.aggregation,
		params.top,
		params.orderby,
		params.filter,
		params.resultType,
		params.metricnamespace,
	)
}

// FetchMetricData returns metric data
func FetchMetricData(ctx context.Context, subscriptionID, resourceGroup, namespace, resource, metricName, aggregation string) (*insights.MetricValue, error) {
	endTime := time.Now().UTC()
	startTime := endTime.Add(time.Duration(-5) * time.Minute)
	timespan := fmt.Sprintf("%s/%s", startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))

	input := &metricsListInput{
		subscriptionID: subscriptionID,
		resourceURI: fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/%s/%s",
			subscriptionID,
			resourceGroup,
			namespace,
			resource,
		),
		timespan:    timespan,
		interval:    to.StringPtr("PT1M"),
		aggregation: aggregation,
		metricnames: metricName,
		resultType:  insights.Data,
	}
	res, err := metricsList(ctx, input)
	if err != nil {
		return nil, err
	}

	var latestData insights.MetricValue
	for _, v := range *res.Value {
		for _, elem := range *v.Timeseries {
			for _, d := range *elem.Data {
				rv := reflect.ValueOf(d)
				av := rv.FieldByName(aggregation)
				if av.IsNil() {
					continue
				}

				if latestData.TimeStamp == nil {
					latestData = d
					continue
				}

				if d.TimeStamp.After(latestData.TimeStamp.Time) {
					latestData = d
				}
			}
		}
	}

	return &latestData, nil
}
