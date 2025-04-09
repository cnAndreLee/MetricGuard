package hwc

import (
	"fmt"
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	ces "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ces/v1"
	cesModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ces/v1/model"
	cesRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ces/v1/region"
	cesV2 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ces/v2"
	cesV2Model "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ces/v2/model"
	cesV2Region "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ces/v2/region"
	cesV3 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ces/v3"
	cesV3Model "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ces/v3/model"
	cesV3Region "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ces/v3/region"
)

var cesClient *ces.CesClient
var cesV2Client *cesV2.CesClient
var cesV3Client *cesV3.CesClient

func cesInit(regionID string, auth *basic.Credentials) {

	// v1
	cesRegion, err := cesRegion.SafeValueOf(regionID)
	if err != nil {
		fmt.Println(err)
		return
	}

	cesHcClient, err := ces.CesClientBuilder().
		WithRegion(cesRegion).
		WithCredential(auth).
		SafeBuild()
	if err != nil {
		fmt.Println(err)
		return
	}
	cesClient = ces.NewCesClient(cesHcClient)

	// v2
	cesV2Region, err := cesV2Region.SafeValueOf(regionID)
	if err != nil {
		fmt.Println(err)
		return
	}

	cesV2HcClient, err := cesV2.CesClientBuilder().
		WithRegion(cesV2Region).
		WithCredential(auth).
		SafeBuild()
	if err != nil {
		fmt.Println(err)
		return
	}
	cesV2Client = cesV2.NewCesClient(cesV2HcClient)

	// v3
	cesV3Region, err := cesV3Region.SafeValueOf(regionID)
	if err != nil {
		fmt.Println(err)
		return
	}

	cesV3HcClient, err := cesV3.CesClientBuilder().
		WithRegion(cesV3Region).
		WithCredential(auth).
		SafeBuild()
	if err != nil {
		fmt.Println(err)
		return
	}
	cesV3Client = cesV3.NewCesClient(cesV3HcClient)
}

// 查询当前可监控的所有指标
func ListMetric() {
	request := &cesModel.ListMetricsRequest{}

	response, err := cesClient.ListMetrics(request)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, metric := range *response.Metrics {
		fmt.Printf("Metric: %v\n", &metric.MetricName)
	}
}

// 查询指标指定服务器id指标
func ShowMetricData(ServerId string) {

	request := &cesModel.ShowMetricDataRequest{
		Namespace:  "SYS.ECS",
		MetricName: "cpu_util",
		From:       time.Now().Add(-5 * time.Minute).UnixMilli(),
		To:         time.Now().UnixMilli(),
		Filter:     cesModel.GetShowMetricDataRequestFilterEnum().AVERAGE,
		Period:     1,
		Dim0:       fmt.Sprintf("instance_id,%s", ServerId),
	}

	response, err := cesClient.ShowMetricData(request)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response)
	for _, v := range *response.Datapoints {
		fmt.Printf("%6f, %v\n", *v.Average, time.UnixMilli(v.Timestamp).Format("2006-01-02 15:04:05"))
	}
}

// 批量查询服务器监控指标
func BatchListMetricData() {
	request := &cesModel.BatchListMetricDataRequest{
		Body: &cesModel.BatchListMetricDataRequestBody{
			Metrics: []cesModel.MetricInfo{
				{
					Namespace:  "AGT.ECS",
					MetricName: "cpu_usage",
					Dimensions: []cesModel.MetricsDimension{
						{
							Name:  "instance_id",
							Value: "08621440-8bec-4dea-b80f-2dfba9aab479",
						},
					},
				},
				{
					Namespace:  "AGT.ECS",
					MetricName: "load_average1",
					Dimensions: []cesModel.MetricsDimension{
						{
							Name:  "instance_id",
							Value: "08621440-8bec-4dea-b80f-2dfba9aab479",
						},
					},
				},
				{
					Namespace:  "AGT.ECS",
					MetricName: "cpu_usage",
					Dimensions: []cesModel.MetricsDimension{
						{
							Name:  "instance_id",
							Value: "501d36e7-4291-45b8-80d6-f5bebf3faa11",
						},
					},
				},
				{
					Namespace:  "AGT.ECS",
					MetricName: "mem_usedPercent",
					Dimensions: []cesModel.MetricsDimension{
						{
							Name:  "instance_id",
							Value: "501d36e7-4291-45b8-80d6-f5bebf3faa11",
						},
					},
				},
			},

			From:   time.Now().Add(-1 * time.Minute).UnixMilli(),
			To:     time.Now().UnixMilli(),
			Filter: "average",
			Period: "1",
		},
	}

	response, err := cesClient.BatchListMetricData(request)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(response)
}

// 获得所有服务器agent状态，返回状态异常的服务器id列表
func ListAgentStatus() []string {

	var listServerId []string

	var listAgentErrorServer []string
	for _, v := range *ServersDetails.Servers {
		listServerId = append(listServerId, v.Id)
	}
	request := &cesV3Model.ListAgentStatusRequest{
		Body: &cesV3Model.ListAgentStatusRequestBody{
			InstanceIds: listServerId,
		},
	}

	response, err := cesV3Client.ListAgentStatus(request)
	if err != nil {
		fmt.Println(err)
		return listAgentErrorServer
	}

	for _, v := range *response.AgentStatus {
		if v.UniagentStatus.Value() != "running" {
			listAgentErrorServer = append(listAgentErrorServer, *v.InstanceId)
		}
		// fmt.Printf("%v,%v \n", *v.InstanceId, v.UniagentStatus.Value())
	}

	return listAgentErrorServer

}

func CesListAlarmRules() {
	limit := int32(100)
	request := &cesV2Model.ListAlarmRulesRequest{
		Limit: &limit,
	}

	response, err := cesV2Client.ListAlarmRules(request)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response)

}
