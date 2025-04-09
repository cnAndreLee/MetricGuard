package hwc

import (
	"fmt"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	iam "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3"
	iamModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"
	iamRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/region"
)

var iamClient *iam.IamClient

func iamInit(regionID string, auth *basic.Credentials) {
	iamRegion, err := iamRegion.SafeValueOf(regionID)
	if err != nil {
		fmt.Println(err)
		return
	}

	iamHcClient, err := iam.IamClientBuilder().
		WithRegion(iamRegion).
		WithCredential(auth).
		SafeBuild()
	if err != nil {
		fmt.Println(err)
		return
	}
	iamClient = iam.NewIamClient(iamHcClient)
}

func ListProjects() {
	request := &iamModel.KeystoneListProjectsRequest{}

	response, err := iamClient.KeystoneListProjects(request)
	if err == nil {
		fmt.Printf("%+v\n", response)
	} else {
		fmt.Println(err)
	}
}
