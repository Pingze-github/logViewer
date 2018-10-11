package main

import (
	"errors"
	"fmt"
	"github.com/Pingze-github/sift"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
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

// 获取指定目录下的文件树
func GetFileTreeHandler(ctx *gin.Context) {
	node, err := getFileTree("G:/raid/youxin.357.com/logs/", 10)
	if err != nil {
		cReturn(ctx, RetBody{Error: err})
		return
	}
	cReturn(ctx, RetBody{Message: gin.H{"root": FileNode2Gin(node)}})
}

// 获取指定文件指定起止行数间的各行内容
func GetFileLinesHandler(ctx *gin.Context) {
	lstartS := ctx.DefaultQuery("lstart", "1")
	lendS := ctx.DefaultQuery("lend", "100")
	lstart, _ := strconv.ParseInt(lstartS, 10, 64)
	lend, _ := strconv.ParseInt(lendS, 10, 64)
	lines, err := getFileLines("G:/raid/youxin.357.com/logs/main/log20180830", lstart, lend + 1)
	if err != nil {
		cReturn(ctx, RetBody{Error: err})
		return
	}
	cReturn(ctx, RetBody{Message: gin.H{"lines": lines2Gin(lines)}})
}

// 搜索指定文件中匹配正则表达式的行
func SearchFileWithPattern(ctx *gin.Context) {
	pattern := ctx.DefaultQuery("pattern", ".*")
	filePath := ctx.DefaultQuery("filePath", "G:/raid/youxin.357.com/logs/main/log20180830")

	pageS := ctx.DefaultQuery("page", "1")
	pagesizeS := ctx.DefaultQuery("pagesize", "100")
	page, _ := strconv.ParseInt(pageS, 10, 64)
	pagesize, _ := strconv.ParseInt(pagesizeS, 10, 64)
	sliceStart := (page - 1) * pagesize
	sliceEnd := page * pagesize

	// 避免搜索.时sift的异常阻塞
	if pattern == "." {
		pattern = ".+"
	}

	cmd := fmt.Sprintf("-e %s %s -n --follow --binary-skip --limit=%d", pattern, filePath, sliceEnd)
	searchResult, err := sift.ExecuteSiftCmd(cmd, time.Duration(15e9))

	if err != nil {
		cReturn(ctx, RetBody{Error: err})
		return
	}

	result := searchResult.Results[0]
	matches := result.Matches
	if int64(len(matches)) < sliceStart {
		sliceStart = int64(len(matches))
	}
	if int64(len(matches)) < sliceEnd {
		sliceEnd = int64(len(matches))
	}
	matches = matches[sliceStart:sliceEnd]

	fmt.Println("执行搜索耗时", searchResult.TimeCost)

	cReturn(ctx, RetBody{Message: gin.H{
		"timeCost": float32(searchResult.TimeCost) / 1000000,
		"results": matches2H(matches),
		"target": result.Target,
	}})
}