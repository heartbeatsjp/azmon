package main

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-09-01/insights"
	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-09-01/insights/insightsapi"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
)

// FetchMetricDataInput is input parameters for FetchMetricData
type FetchMetricDataInput struct {
	subscriptionID string
	resourceGroup  string
	namespace      string
	resource       string
	metricNames    []string
	aggregation    string
	startTime      time.Time
	endTime        time.Time
	interval       int
}

// FetchMetricDefinitionsInput is input parameters for FetchMetricDefinitions
type FetchMetricDefinitionsInput struct {
	subscriptionID  string
	resourceGroup   string
	namespace       string
	resource        string
	resourceURI     string
	metricnamespace string
	startTime       time.Time
	endTime         time.Time
	interval        int
}

// Client is an API Client for Azure
type Client struct {
	MetricsClient           insightsapi.MetricsClientAPI
	MetricDefinitionsClient insightsapi.MetricDefinitionsClientAPI
}

// NewClient returns *Client with setting Authorizer
func NewClient(subscriptionID string) (*Client, error) {
	a, err := auth.NewAuthorizerFromFile(azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		return &Client{}, err
	}

	metricsClient := insights.NewMetricsClient(subscriptionID)
	metricsClient.Authorizer = a

	metricDefinitionsClient := insights.NewMetricDefinitionsClient(subscriptionID)
	metricDefinitionsClient.Authorizer = a

	return &Client{
		MetricsClient:           metricsClient,
		MetricDefinitionsClient: metricDefinitionsClient,
	}, nil
}

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

func (c *Client) metricDefinitionsList(ctx context.Context, params *metricDefinitionsListInput) (insights.MetricDefinitionCollection, error) {
	return c.MetricDefinitionsClient.List(
		ctx,
		params.resourceURI,
		params.metricnamespace,
	)
}

func (c *Client) metricsList(ctx context.Context, params *metricsListInput) (insights.Response, error) {
	return c.MetricsClient.List(
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

// FetchMetricDefinitions returns metric definitions
func FetchMetricDefinitions(ctx context.Context, c *Client, params FetchMetricDefinitionsInput) (*[]insights.MetricDefinition, error) {
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
	res, err := c.metricDefinitionsList(ctx, input)
	if err != nil {
		return nil, err
	}
	return res.Value, nil
}

// FetchMetricData fetches metric data and returns latest value with metric name as hash key
func FetchMetricData(ctx context.Context, c *Client, params FetchMetricDataInput) (map[string]*insights.MetricValue, error) {
	timespan := fmt.Sprintf("%s/%s", params.startTime.Format(time.RFC3339), params.endTime.Format(time.RFC3339))

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

	metrics := make(map[string]*insights.MetricValue)
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
			interval:    to.StringPtr(fmt.Sprintf("PT%dS", params.interval)),
			aggregation: params.aggregation,
			metricnames: m,
			resultType:  insights.Data,
		}
		res, err := c.metricsList(ctx, input)
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

				if latestData.TimeStamp != nil {
					metrics[*v.Name.Value] = &latestData
				}
			}
		}
	}
	return metrics, nil
}
