package controller

import (
	"fmt"
	"net/http"

	"github.com/cnAndreLee/MetricGuard/hwc"
	"github.com/cnAndreLee/MetricGuard/model"
	"github.com/cnAndreLee/MetricGuard/response"

	"github.com/gin-gonic/gin"
)

func Alert(c *gin.Context) {

	// 接收客户端DTO
	var request model.AltreRequest
	err := c.ShouldBind(&request)
	if err != nil {
		res := response.ResponseStruct{
			HttpStatus: http.StatusOK,
			Code:       response.FailCode,
			Msg:        "数据无效",
			Data:       nil,
		}

		response.Response(c, res)
		return
	}

	if request.Host == "" || request.Ip == "" || request.Type == "" {
		res := response.ResponseStruct{
			HttpStatus: http.StatusOK,
			Code:       response.FailCode,
			Msg:        "数据无效",
			Data:       nil,
		}

		response.Response(c, res)
		return
	}

	msg := ""
	levelString := ""
	if request.Type == "5" {
		levelString = "五级告警"
	} else if request.Type == "4" {
		levelString = "四级告警"
	} else if request.Type == "0" {
		levelString = "告警恢复"
	} else {
		res := response.ResponseStruct{
			HttpStatus: http.StatusOK,
			Code:       response.FailCode,
			Msg:        "数据无效",
			Data:       nil,
		}

		response.Response(c, res)
		return
	}

	msg = fmt.Sprintf("[%v]\n主机:%v\nIP:%v\n%v(%v):%v",
		levelString,
		request.Host,
		request.Ip,
		request.MetricCnName,
		request.MetricUnit,
		request.MetricValue,
	)

	// fmt.Println(request)
	// fmt.Println(msg)

	if request.Type == "5" {
		err = hwc.PublishMessage(msg, "5")
	} else if request.Type == "4" {
		err = hwc.PublishMessage(msg, "4")
	} else if request.Type == "0" {
		if request.LastType == "4" {
			err = hwc.PublishMessage(msg, "4")
		} else if request.LastType == "5" {
			err = hwc.PublishMessage(msg, "5")
		}
	}

	if err != nil {
		res := response.ResponseStruct{
			HttpStatus: http.StatusOK,
			Code:       response.FailCode,
			Msg:        err.Error(),
			Data:       nil,
		}

		response.Response(c, res)
		return
	}

	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Msg:        "",
		Data:       nil,
	}

	response.Response(c, res)

}
