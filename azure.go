package main

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-11-01-preview/insights"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
)

type metricDefinitionsListInput struct {
	subscriptionID  string
	resourceURI     string
	metricnamespace string
}

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

func metricDefinitionsList(ctx context.Context, params *metricDefinitionsListInput) (insights.MetricDefinitionCollection, error) {
	client := insights.NewMetricDefinitionsClient(params.subscriptionID)
	a, err := auth.NewAuthorizerFromFile(azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		return insights.MetricDefinitionCollection{}, err
	}
	client.Authorizer = a
	return client.List(
		ctx,
		params.resourceURI,
		params.metricnamespace,
	)
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

func FetchMetricDefinitions(ctx context.Context, params FetchMetricDefinitionsInput) (*[]insights.MetricDefinition, error) {
	input := &metricDefinitionsListInput{
		subscriptionID: params.subscriptionID,
		resourceURI: fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/%s/%s",
			params.subscriptionID,
			params.resourceGroup,
			params.namespace,
			params.resource,
		),
		metricnamespace: params.metricnamespace,
	}
	res, err := metricDefinitionsList(ctx, input)
	if err != nil {
		return nil, err
	}
	return res.Value, nil
}

// FetchMetricData returns metric data
func FetchMetricData(ctx context.Context, params FetchMetricDataInput) (map[string]insights.MetricValue, error) {
	endTime := time.Now().UTC()
	startTime := endTime.Add(time.Duration(-5) * time.Minute)
	timespan := fmt.Sprintf("%s/%s", startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))

	var metricNames []string
	const metricsCountLimitPerRequest int = 20
	for {
		if len(params.metricNames) <= metricsCountLimitPerRequest {
			metricNames = append(metricNames, strings.Join(params.metricNames, ","))
			break
		}

		metricNames = append(metricNames, strings.Join(params.metricNames[:metricsCountLimitPerRequest], ","))
		params.metricNames = params.metricNames[metricsCountLimitPerRequest:]
	}

	metrics := make(map[string]insights.MetricValue)
	for _, m := range metricNames {
		input := &metricsListInput{
			subscriptionID: params.subscriptionID,
			resourceURI: fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/%s/%s",
				params.subscriptionID,
				params.resourceGroup,
				params.namespace,
				params.resource,
			),
			timespan:    timespan,
			interval:    to.StringPtr("PT1M"),
			aggregation: params.aggregation,
			metricnames: m,
			resultType:  insights.Data,
		}
		res, err := metricsList(ctx, input)
		if err != nil {
			return nil, err
		}

		for _, v := range *res.Value {
			for _, elem := range *v.Timeseries {
				var latestData insights.MetricValue
				for _, d := range *elem.Data {
					rv := reflect.ValueOf(d)
					av := rv.FieldByName(params.aggregation)
					if av.IsNil() {
						continue
					}

					if d.TimeStamp == nil {
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
				//for debug
				//fmt.Printf("%s: %v\n", *v.Name.Value, latestData)
				if latestData.TimeStamp != nil {
					metrics[*v.Name.Value] = latestData
				}
			}
		}
	}
	return metrics, nil
}
