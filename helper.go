package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RetBody struct {
	Succeed bool
	Status int
	String string
	Html string
	Raw []byte
	ContentType string
	Code int
	Desc string
	Message interface{}
	Error error
}

type JSONRet struct {
	Succeed bool `json:"Succeed"`
	Code int `json:"Code"`
	Desc string `json:"Desc"`
	Message interface{}
}

// 上下文统一返回
func cReturn(ctx *gin.Context, ret RetBody) {
	// 按条件返回
	// 优先级 Error > String > Html > Raw > Others
	if ret.Error != nil {
		cReturnStructError(ctx, ret)
	} else if ret.String != "" {
		ctx.Data(ret.Status, "text/plain", []byte(ret.String))
	} else if ret.Html != "" {
		ctx.Data(ret.Status, "text/html", []byte(ret.Html))
	} else if ret.Raw != nil {
		if ret.ContentType == "" {
			ret.ContentType = "text/plain"
		}
		ctx.Data(ret.Status, ret.ContentType, []byte(ret.Raw))
	} else {
		// 结构化返回
		cReturnStruct(ctx, ret)
	}
	ctx.Abort()
}

// 上下文结构化返回
func cReturnStruct(ctx *gin.Context, ret RetBody)  {
	// 设置默认值
	if ret.Code == 0 {
		ret.Code = http.StatusOK
	}
	if ret.Desc == "" {
		ret.Desc = "成功"
	}
	if ret.Message == nil {
		ret.Message = gin.H{}
	}
	if ret.Status == 0 && 100 <= ret.Code && 600 >= ret.Code{
		ret.Status = ret.Code
	}

	ret.Succeed = ret.Code == http.StatusOK

	jsonRet := JSONRet{Succeed: ret.Succeed, Code: ret.Code, Desc: ret.Desc, Message: ret.Message}

	ctx.JSON(ret.Status, jsonRet)
}

// 错误返回
func cReturnStructError(ctx *gin.Context, ret RetBody)  {
	errorMessage := fmt.Sprint(ret.Error)
	if errorMessage == "" {
		errorMessage = "内部错误"
	}
	jsonRet := JSONRet{Succeed: false, Code: http.StatusInternalServerError, Desc: errorMessage, Message: gin.H{}}
	ctx.JSON(http.StatusInternalServerError, jsonRet)
}