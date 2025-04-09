package hwc

import (
	"errors"
	"fmt"
	"log"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	smn "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/smn/v2"
	smnModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/smn/v2/model"
	smnRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/smn/v2/region"
)

var smnClient *smn.SmnClient

func smnInit(regionID string, auth *basic.Credentials) {
	region, err := smnRegion.SafeValueOf(regionID)
	if err != nil {
		fmt.Println(err)
		return
	}

	hcClient, err := smn.SmnClientBuilder().
		WithRegion(region).
		WithCredential(auth).
		SafeBuild()
	if err != nil {
		fmt.Println(err)
		return
	}
	smnClient = smn.NewSmnClient(hcClient)
}

// 添加订阅者
func AddSubscription(phoneNumber string, remark string, topic string) {
	request := &smnModel.AddSubscriptionRequest{}
	if topic == "4" {
		request.TopicUrn = topicurn4
	} else if topic == "5" {
		request.TopicUrn = topicurn5
	} else {
		log.Println("wrong topic")
		return
	}

	request.Body = &smnModel.AddSubscriptionRequestBody{
		Protocol: "sms",
		Endpoint: phoneNumber,
		Remark:   &remark,
	}

	response, err := smnClient.AddSubscription(request)
	if err == nil {
		fmt.Printf("%+v\n", response)
	} else {
		fmt.Println(err)
	}
}

func PublishMessage(message string, level string) error {

	request := &smnModel.PublishMessageRequest{}

	if level == "5" {
		request.TopicUrn = topicurn5
	} else if level == "4" {
		request.TopicUrn = topicurn4
	} else {
		log.Println("PublishMessage: wrong topic")
		return errors.New("wrong topic")
	}

	messagePublishMessageRequestBody := message
	request.Body = &smnModel.PublishMessageRequestBody{
		// TimeToLive:        &timeToLivePublishMessageRequestBody,
		Message: &messagePublishMessageRequestBody,
		//Subject:           &subjectPublishMessageRequestBody,
	}

	_, err := smnClient.PublishMessage(request)
	if err != nil {
		return err
	}

	return nil
}

// func listSubscriptions(auth *basic.Credentials, regionID string) (*smnModel.ListSubscriptionsResponse, error) {
func ListSubscriptions() {

	valueRequestStatus := int32(1)
	valueRequestProtocol := "sms"
	request := &smnModel.ListSubscriptionsRequest{
		Protocol: &valueRequestProtocol,
		Status:   &valueRequestStatus,
	}

	response, err := smnClient.ListSubscriptions(request)
	if err == nil {
		fmt.Printf("%+v\n", response)
	} else {
		fmt.Println(err)
	}
	// return response, err
}

// 创建主题
func CreateTopic(TopicName string, TopicDisplayName string, ProjectId string) error {
	request := &smnModel.CreateTopicRequest{}

	request.Body = &smnModel.CreateTopicRequestBody{
		Name:                TopicName,
		DisplayName:         TopicDisplayName,
		EnterpriseProjectId: &ProjectId,
	}

	response, err := smnClient.CreateTopic(request)
	if err == nil {
		fmt.Printf("%+v\n", response)
		return err
	} else {
		fmt.Println(err)
	}

	return nil
}

// 查询主题详情
func ListTopicDetails(Topic string) {

	request := &smnModel.ListTopicDetailsRequest{
		TopicUrn: Topic,
	}

	response, err := smnClient.ListTopicDetails(request)
	if err == nil {
		fmt.Printf("%+v\n", response)
	} else {
		fmt.Println(err)
	}
}
