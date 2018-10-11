package main

import (
	"bufio"
	"errors"
	"github.com/Pingze-github/sift"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// 文件行
type Line struct {
	// 行号
	lineno int64
	// 内容
	content string
}

// 执行sift命令（带超时限制）
// timeout为毫秒
//func ExecuteSiftCmdWithTimeout(cmd string, timeout int64) (sift.SearchResult, error) {
//	type SearchResultWithError struct {
//		SearchResult sift.SearchResult
//		err error
//	}
//
//	ch := make(chan SearchResultWithError)
//
//	go func() {
//		searchResult, err := sift.ExecuteSiftCmd(cmd, timeout)
//		ch <- SearchResultWithError{searchResult, err}
//	}()
//
//	select {
//	case resultWithError := <-ch:
//		return resultWithError.SearchResult, resultWithError.err
//	case <-time.After(time.Millisecond * time.Duration(timeout)):
//		return sift.SearchResult{}, errors.New(fmt.Sprintf("sift search timeout for %dms limit", timeout))
//	}
//}

func results2H(results []*sift.Result) []gin.H {
	var resultsH []gin.H
	for _, result := range(results) {
		resultsH = append(resultsH, gin.H{
			"target": result.Target,
			"matches": matches2H(result.Matches),
		})
	}
	return resultsH
}

func matches2H(matches sift.Matches) []gin.H {
	var resultsH []gin.H
	for _, match := range(matches) {
		resultsH = append(resultsH, gin.H{
			"start": match.Start,
			"end": match.End,
			"lineStart": match.LineStart,
			"lineEnd": match.LineEnd,
			"match": match.Match,
			"line": match.Line,
			"lineno": match.Lineno,
		})
	}
	return resultsH
}

// 获取文件的指定函数内容
func getFileLines(filePath string, startLineno int64, endLineno int64) ([]*Line, error) {
	var lines []*Line

	if startLineno > endLineno {
		return lines, errors.New("startLineno must litter than endLineno")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return lines, err;
	}

	buf := bufio.NewReader(file)
	lineno := int64(1)
	for {
		line, err := buf.ReadString('\n')
		if lineno == endLineno {
			break
		}
		if lineno >= startLineno {
			line = strings.TrimSpace(line)
			lines = append(lines, &Line{lineno, line})
		}
		if err != nil {
			if err == io.EOF {
				return lines, nil
			}
			return lines, err
		}
		lineno ++
	}

	if err != nil {
		return lines, err;
	}
	return lines, nil
}

// 将文件行转化为gin.H数组
func lines2Gin(lines []*Line) []gin.H {
	var data []gin.H
	for _, linePoint := range lines{
		line := *linePoint
		lineH := gin.H{
			"lineno": line.lineno,
			"content": line.content,
		}
		data = append(data, lineH)
	}
	return data
}

// 文件系统树节点
type FileNode struct {
	// 文件名
	name string
	// 文件路径
	path string
	// 是否为文件（否则为目录）（win中快捷方式视为文件）
	isFile bool
	// 子节点
	children []*FileNode
}

// 递归获取文件系统树
// TODO 支持linux软链接
// 由于os.Readlink不支持win下快捷方式，不进行特别支持
func getFileTree(dirPath string, limit int64) (FileNode, error) {
	var nodeChildren []*FileNode
	node := FileNode{name:dirPath, path:dirPath, isFile:false}
	rd, err := ioutil.ReadDir(dirPath)
	for _, fi := range rd {
		fiName := fi.Name()
		fiPath := dirPath + "/" + fiName
		if fi.IsDir() {
			fiNode, err := getFileTree(fiPath, limit)
			if err != nil {
				return node, err
			}
			nodeChildren = append(nodeChildren, &fiNode)
		} else {
			fiNode := FileNode{name:fiName, path:fiPath, isFile:true}
			nodeChildren = append(nodeChildren, &fiNode)
		}
	}
	node.children = nodeChildren
	return node, err
}

// 将文件树转化为gin.H格式
func FileNode2Gin(node FileNode) gin.H {
	var nodeChildren []gin.H
	for _, nodeChildPoint := range node.children{
		nodeChild := *nodeChildPoint
		ginNodeChild := gin.H{
			"name": nodeChild.name,
			"path": nodeChild.path,
			"isFile": nodeChild.isFile,
			"chidren": FileNode2Gin(nodeChild),
		}
		nodeChildren = append(nodeChildren, ginNodeChild)
	}

	return gin.H{
		"name": node.name,
		"path": node.path,
		"isFile": node.isFile,
		"chidren": nodeChildren,
	}
}
