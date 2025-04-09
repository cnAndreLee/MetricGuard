package config

import (
	"errors"
	"fmt"
	"os"
)

type ConfigStruct struct {
	HW HWConfig
	HS HttpServer
}

type HWConfig struct {
	Ak        string
	Sk        string
	ProjectId string
	RegionId  string
	TopicUrn5 string
	TopicUrn4 string
}

type HttpServer struct {
	Port string
}

var CONFIG ConfigStruct

func Init() (*ConfigStruct, error) {

	EnvMap := map[string]string{
		"HWSMN_HW_AK":        "",
		"HWSMN_HW_SK":        "",
		"HWSMN_HW_PROJECTID": "",
		"HWSMN_HW_REGIONID":  "",
		"HWSMN_HW_TOPICURN5": "",
		"HWSMN_HW_TOPICURN4": "",
		"HWSMN_HTTP_PORT":    "",
	}

	errExists := false
	for i := range EnvMap {
		if str, exists := os.LookupEnv(i); !exists {
			fmt.Println("env: ", i, "not found")
			errExists = true
		} else {
			EnvMap[i] = str
		}
	}
	if errExists {
		return nil, errors.New("env check error")
	}

	CONFIG = ConfigStruct{
		HW: HWConfig{
			Ak:        EnvMap["HWSMN_HW_AK"],
			Sk:        EnvMap["HWSMN_HW_SK"],
			RegionId:  EnvMap["HWSMN_HW_REGIONID"],
			ProjectId: EnvMap["HWSMN_HW_PROJECTID"],
			TopicUrn5: EnvMap["HWSMN_HW_TOPICURN5"],
			TopicUrn4: EnvMap["HWSMN_HW_TOPICURN4"],
		},
		HS: HttpServer{
			Port: EnvMap["HWSMN_HTTP_PORT"],
		},
	}

	fmt.Printf("Config loaded --- %+v \n", CONFIG)
	return &CONFIG, nil
}
