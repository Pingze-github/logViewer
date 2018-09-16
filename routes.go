package main

import "github.com/gin-gonic/gin"

func handleRoutes(router *gin.Engine) {
	router.GET("/", IndexHandler)
	router.GET("/api", ApiHandler)
	router.GET("/error", ErrorHandler)
}
