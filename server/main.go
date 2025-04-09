package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/cnAndreLee/MetricGuard/config"
	"github.com/cnAndreLee/MetricGuard/hwc"
	"github.com/cnAndreLee/MetricGuard/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	_, err := config.Init()
	if err != nil {
		fmt.Println(err)
		return
	}

	hwc.Init()

	// smn
	// hwc.ListSubscriptions()
	// hwc.AddSubscription("+8615596099788", "李帅", "5")
	// hwc.ListTopicDetails("urn:smn:cn-north-4:3fb032df961740068fbc86c7e6f1058d:app_test_demo")
	// hwc.CreateTopic("jsz_alert_5", "应用五级告警", "0")

	// hwc.CesListAlarmRules()

	// 初始化结构体，重要
	// hwc.ListServersDetails()
	// hwc.PrintServers()

	// hwc.BatchListMetricData()
	// hwc.ShowMetricData("08621440-8bec-4dea-b80f-2dfba9aab479")

	// fmt.Println("Agent状态异常的服务器:")
	// hwc.PrintServerInfo(hwc.ListAgentStatus())
	// fmt.Println("---------------------")

	r := gin.Default()

	r = routes.CollectRoute(r)
	// port := myconfig.HS.Port
	// fmt.Println(port)

	socketPath := "/tmp/metricguard.sock"
	// 确保旧的 socket 文件被删除
	if err := os.RemoveAll(socketPath); err != nil {
		log.Fatalf("Failed to remove old socket: %v", err)
	}

	// 创建 Unix listener
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatalf("Failed to listen on Unix socket: %v", err)
	}

	// 设置 socket 文件权限 (可选)
	if err := os.Chmod(socketPath, 0777); err != nil {
		log.Printf("Warning: Failed to change socket permissions: %v", err)
	}

	// 优雅关闭处理
	setupGracefulShutdown(listener, socketPath)

	// 启动服务
	log.Printf("Server is listening on Unix socket: %s", socketPath)
	if err := r.RunListener(listener); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

	// 抽象unix socket
	// listener, err := net.Listen("unix", "@metric_guard")
	// if err != nil {
	// 	panic(err)
	// }
	// panic(r.RunListener(listener))
	// panic(r.Run(":" + port))

}

func setupGracefulShutdown(listener net.Listener, socketPath string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down server...")
		listener.Close()
		os.Remove(socketPath)
		os.Exit(0)
	}()
}
