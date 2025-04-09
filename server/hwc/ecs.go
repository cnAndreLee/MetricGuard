package hwc

import (
	"fmt"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	ecs "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2"
	ecsModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/model"
	ecsRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/region"
)

var ecsClient *ecs.EcsClient

func ecsInit(regionID string, auth *basic.Credentials) {
	ecsRegion, err := ecsRegion.SafeValueOf(regionID)
	if err != nil {
		fmt.Println(err)
		return
	}

	ecsHcClient, err := ecs.EcsClientBuilder().
		WithRegion(ecsRegion).
		WithCredential(auth).
		SafeBuild()
	if err != nil {
		fmt.Println(err)
		return
	}
	ecsClient = ecs.NewEcsClient(ecsHcClient)
}

// 列出所有服务器信息
// https://support.huaweicloud.com/api-ecs/zh-cn_topic_0094148850.html
func ListServersDetails() {
	limit := int32(1000)
	request := &ecsModel.ListServersDetailsRequest{
		Limit: &limit,
	}

	response, err := ecsClient.ListServersDetails(request)
	if err != nil {
		return
	}
	ServersDetails = response

	// fmt.Println(*response)
	// for _, v := range *response.Servers {
	// 	vpc_id := v.Metadata["vpc_id"]
	// 	AllServers = append(AllServers, Server{
	// 		Id:      v.Id,
	// 		Name:    v.Name,
	// 		Address: v.Addresses[vpc_id][0].Addr,
	// 	})
	// }
}
