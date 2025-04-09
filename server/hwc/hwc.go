package hwc

import (
	"fmt"

	"github.com/cnAndreLee/MetricGuard/config"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	ecsModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/model"
)

var ServersDetails *ecsModel.ListServersDetailsResponse

var topicurn5 string
var topicurn4 string

func Init() {

	ak := config.CONFIG.HW.Ak
	sk := config.CONFIG.HW.Sk
	regionID := config.CONFIG.HW.RegionId
	projectId := config.CONFIG.HW.ProjectId
	topicurn5 = config.CONFIG.HW.TopicUrn5
	topicurn4 = config.CONFIG.HW.TopicUrn4

	auth, err := basic.NewCredentialsBuilder().
		WithAk(ak).
		WithSk(sk).
		WithProjectId(projectId).
		SafeBuild()
	if err != nil {
		fmt.Println(err)
		return
	}

	smnInit(regionID, auth)
	ecsInit(regionID, auth)
	iamInit(regionID, auth)
	cesInit(regionID, auth)
	cbhInit(regionID, auth)
}

func PrintServers() {
	for _, v := range *ServersDetails.Servers {
		fmt.Printf("%v %v %v %v\n", v.Id, v.Name, v.Addresses[v.Metadata["vpc_id"]], v.Metadata["os_type"])
	}
}

func PrintServerInfo(listServer []string) {
	for _, v := range listServer {
		for _, v2 := range *ServersDetails.Servers {
			if v == v2.Id {
				fmt.Printf("%v,%v,%v,%v\n", v, v2.Name, v2.Addresses[v2.Metadata["vpc_id"]][0].Addr, v2.Metadata["os_type"])
			}
		}
	}
}
