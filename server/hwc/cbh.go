package hwc

import (
	"fmt"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	cbh "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cbh/v2"
	cbhModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cbh/v2/model"
	cbhRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cbh/v2/region"
)

var cbhClient *cbh.CbhClient

func cbhInit(regionID string, auth *basic.Credentials) {
	cbhRegion, err := cbhRegion.SafeValueOf(regionID)
	if err != nil {
		fmt.Println(err)
		return
	}

	cbhHcClient, err := cbh.CbhClientBuilder().
		WithRegion(cbhRegion).
		WithCredential(auth).
		SafeBuild()
	if err != nil {
		fmt.Println(err)
		return
	}
	cbhClient = cbh.NewCbhClient(cbhHcClient)
}

// 查询堡垒机实例列表
func ListCbhInstances() {
	request := &cbhModel.ListInstancesRequest{}

	response, err := cbhClient.ListInstances(request)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response)
}
