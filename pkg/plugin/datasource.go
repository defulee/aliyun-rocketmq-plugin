package plugin

import (
	"context"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	ons "github.com/aliyun/alibaba-cloud-sdk-go/services/ons"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-starter-datasource-backend/pkg/models"
)

// Make sure RocketMqDatasource implements required interfaces. This is important to do
// since otherwise we will only get a not implemented error response from plugin in
// runtime. In this example datasource instance implements backend.QueryDataHandler,
// backend.CheckHealthHandler, backend.StreamHandler interfaces. Plugin should not
// implement all these interfaces - only those which are required for a particular task.
// For example if plugin does not need streaming functionality then you are free to remove
// methods that implement backend.StreamHandler. Implementing instancemgmt.InstanceDisposer
// is useful to clean up resources used by previous datasource instance when a new datasource
// instance created upon datasource Settings changed.
var (
	_ backend.QueryDataHandler      = (*RocketMqDatasource)(nil)
	_ backend.CheckHealthHandler    = (*RocketMqDatasource)(nil)
	_ instancemgmt.InstanceDisposer = (*RocketMqDatasource)(nil)
)

// RocketMqDatasource is a datasource which can respond to data queries, reports its health.
type RocketMqDatasource struct {
	Client   *ons.Client
	Settings *models.PluginSettings
	log      log.Logger
}

// NewRocketMqDatasource creates a new datasource instance.
func NewRocketMqDatasource(settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	log.DefaultLogger.Info("NewRocketMqDatasource called")
	pluginSettings, _ := models.LoadPluginSettings(settings)
	log.DefaultLogger.Info("NewRocketMqDatasource pluginSettings",
		"AccessKeyId", pluginSettings.AccessKeyId,
		"AccessKeySecret", pluginSettings.Secrets.AccessKeySecret,
		"Region", pluginSettings.Region,
		"InstanceId", pluginSettings.InstanceId)
	config := sdk.NewConfig()

	credential := credentials.NewAccessKeyCredential(pluginSettings.AccessKeyId, pluginSettings.Secrets.AccessKeySecret)
	/* use STS Token
	credential := credentials.NewStsTokenCredential("<your-access-key-id>", "<your-access-key-secret>", "<your-sts-token>")
	*/
	client, err := ons.NewClientWithOptions(pluginSettings.Region, config, credential)
	if err != nil {
		panic(err)
	}

	return &RocketMqDatasource{
		Client:   client,
		Settings: pluginSettings,
		log:      log.DefaultLogger,
	}, nil
}

// Dispose here tells plugin SDK that plugin wants to clean up resources when a new instance
// created. As soon as datasource Settings change detected by SDK old datasource instance will
// be disposed and a new one will be created using NewRocketMqDatasource factory function.
func (d *RocketMqDatasource) Dispose() {
	// Clean up datasource instance resources.
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifier).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (d *RocketMqDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	d.log.Info("RocketMqDatasource QueryData called", "request", req)

	// create response struct
	response := backend.NewQueryDataResponse()

	// loop over queries and execute them individually.
	for _, q := range req.Queries {
		res := d.query(ctx, req.PluginContext, q)

		// save the response in a hashmap
		// based on with RefID as identifier
		response.Responses[q.RefID] = res
	}

	return response, nil
}

func (d *RocketMqDatasource) query(_ context.Context, pCtx backend.PluginContext, query backend.DataQuery) backend.DataResponse {
	response := backend.DataResponse{}

	payload, err := models.ParsePayload(query)
	if err != nil || payload.Hide {
		return response
	}

	// create data frame response.
	// add fields.
	// add the frames to the response.
	if payload.Action == "ConsumerAccumulate" {
		d.queryConsumerAccumulate(payload, query, &response)
	} else if payload.Action == "TrendTopicInputTps" {
		d.queryTrendTopicInputTps(payload, query, &response)
	} else if payload.Action == "TrendGroupOutputTps" {
		d.queryTrendGroupOutputTps(payload, query, &response)
	}

	return response
}

func (d *RocketMqDatasource) queryConsumerAccumulate(payload *models.QueryPayload, query backend.DataQuery, response *backend.DataResponse) {
	request := ons.CreateOnsConsumerAccumulateRequest()

	request.Scheme = "https"

	request.InstanceId = d.Settings.InstanceId
	request.Detail = requests.NewBoolean(true)
	request.GroupId = payload.GroupId

	onsResp, err := d.Client.OnsConsumerAccumulate(request)
	if err != nil {
		fmt.Print(err.Error())
		d.log.Error("OnsConsumerAccumulate ", "GroupId", payload.GroupId, "error ", err)
		return
	}

	d.log.Info("query OnsConsumerAccumulate ", "topic count", len(onsResp.Data.DetailInTopicList.DetailInTopicDo))

	frame := data.NewFrame(query.RefID)

	fieldValArrMap := make(map[string]float64)
	for _, detailInTopicDo := range onsResp.Data.DetailInTopicList.DetailInTopicDo {
		fieldValArrMap[detailInTopicDo.Topic] = float64(detailInTopicDo.TotalDiff)
	}

	for field, val := range fieldValArrMap {
		frame.Fields = append(frame.Fields, data.NewField(field, nil, val))
	}

	response.Frames = append(response.Frames, frame)
}

func (d *RocketMqDatasource) queryTrendTopicInputTps(payload *models.QueryPayload, query backend.DataQuery, response *backend.DataResponse) {
	req := ons.CreateOnsTrendTopicInputTpsRequest()

	req.Scheme = "https"
	req.InstanceId = d.Settings.InstanceId
	req.Topic = payload.Topic
	req.Type = requests.NewInteger(1)
	req.Period = requests.NewInteger(1)
	req.BeginTime = requests.NewInteger(payload.From)
	req.EndTime = requests.NewInteger(payload.To)

	resp, err := d.Client.OnsTrendTopicInputTps(req)
	if err != nil {
		fmt.Print(err.Error())
		d.log.Error("queryTrendTopicInputTps ", "Topic", payload.GroupId, "error ", err)
		return
	}

	d.log.Info("query OnsTrendTopicInputTps ", "record count", len(resp.Data.Records.StatsDataDo))

	frame := data.NewFrame(query.RefID)

	fieldValArrMap := make(map[string]float64)
	for _, statsDataDo := range resp.Data.Records.StatsDataDo {
		fieldValArrMap["time"] = float64(statsDataDo.X)
		fieldValArrMap["value"] = float64(statsDataDo.Y)
	}

	for field, val := range fieldValArrMap {
		frame.Fields = append(frame.Fields, data.NewField(field, nil, val))
	}

	response.Frames = append(response.Frames, frame)
}

func (d *RocketMqDatasource) queryTrendGroupOutputTps(payload *models.QueryPayload, query backend.DataQuery, response *backend.DataResponse) {
	request := ons.CreateOnsTrendGroupOutputTpsRequest()

	request.Scheme = "https"

	request.InstanceId = d.Settings.InstanceId
	request.GroupId = payload.GroupId
	request.Topic = payload.Topic
	request.Type = requests.NewInteger(1)
	request.Period = requests.NewInteger(1)
	request.BeginTime = requests.NewInteger(payload.From)
	request.EndTime = requests.NewInteger(payload.To)

	resp, err := d.Client.OnsTrendGroupOutputTps(request)
	if err != nil {
		fmt.Print(err.Error())
		d.log.Error("OnsTrendGroupOutputTps ", "GroupId", payload.GroupId, "error ", err)
		return
	}

	d.log.Info("query OnsTrendGroupOutputTps ", "topic count", len(resp.Data.Records.StatsDataDo))

	frame := data.NewFrame(query.RefID)

	fieldValArrMap := make(map[string]float64)
	for _, statsDataDo := range resp.Data.Records.StatsDataDo {
		fieldValArrMap["time"] = float64(statsDataDo.X)
		fieldValArrMap["value"] = float64(statsDataDo.Y)
	}

	for field, val := range fieldValArrMap {
		frame.Fields = append(frame.Fields, data.NewField(field, nil, val))
	}

	response.Frames = append(response.Frames, frame)
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (d *RocketMqDatasource) CheckHealth(_ context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	d.log.Info("CheckHealth called", "request", req)

	request := ons.OnsInstanceBaseInfoRequest{
		InstanceId: d.Settings.InstanceId,
	}
	request.Scheme = "https"

	_, err := d.Client.OnsInstanceBaseInfo(&request)
	if err != nil {
		d.log.Info("CheckHealth failed", "error", err)
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: "GetLogStore error",
		}, nil
	}

	d.log.Info("CheckHealth success")
	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "Data source is working",
	}, nil
}
