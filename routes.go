package main

import "github.com/gin-gonic/gin"

func handleRoutes(router *gin.Engine) {
	router.GET("/", IndexHandler)
	router.GET("/api", ApiHandler)
	router.GET("/error", ErrorHandler)

	router.GET("/api/filetree", GetFileTreeHandler)

	router.GET("/api/lines", GetFileLinesHandler)

	router.GET("/api/search", SearchFileWithPattern)
}
