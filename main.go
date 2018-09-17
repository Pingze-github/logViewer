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
	getDirTree("G:/raid/youxin.357.com/logs")
	// server(gin.ReleaseMode, 10001)
}