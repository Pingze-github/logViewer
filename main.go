package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func server(mode string, port int64) {
	gin.SetMode(mode)

	//获得路由实例
	router := gin.Default()

	// 静态资源
	// router.StaticFile("/", "./public/index.html")

	// 注册接口
	handleRoutes(router)

	// 启动服务
	fmt.Println(`Server start @`, port)
	router.Run(fmt.Sprintf(":%d", port))
}

func main() {
	//pattern := "a"
	//
	//if pattern == "." {
	//	pattern = ".*"
	//}
	//
	//filePath := "G:/"
	//// filePath = "tools.go"
	//cmd := fmt.Sprintf("-e %s %s -n --follow --binary-skip --cores=2", pattern, filePath)

	//cmd := "-e ERRO G:/raid/youxin.357.com/logs/main/log20180830 -n --follow --binary-skip --limit=100"
	//
	//fmt.Println("sift " + cmd)
	//
	//fmt.Println("开始执行")
	//searchResult, err := sift.ExecuteSiftCmd(cmd, time.Duration(1e8))
	//if err != nil {
	//	fmt.Println("超时", err.Error())
	//}
	//fmt.Println("0", searchResult)

	//fmt.Println("开始执行")
	//searchResult1, err1 := sift.ExecuteSiftCmd(cmd, time.Duration(1e9))
	//if err1 != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println("0", searchResult1)


	server("debug", 6600)
}