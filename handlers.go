package main

import (
	"errors"
	"github.com/gin-gonic/gin"
)

func IndexHandler(ctx *gin.Context) {
	cReturn(ctx, RetBody{String: "logViewer"})
}

func ApiHandler(ctx *gin.Context) {
	cReturn(ctx, RetBody{})
}

func ErrorHandler(ctx *gin.Context) {
	cReturn(ctx, RetBody{Error: errors.New("")})
}