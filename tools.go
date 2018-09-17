package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/Pingze-github/sift"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Line struct {
	lineno int64
	content string
}

// 执行sift命令（带超时限制）
func ExecuteSiftCmdWithTimeout(cmd string, timeout int64) (sift.SearchResult, error) {

	type SearchResultWithError struct {
		SearchResult sift.SearchResult
		err error
	}

	ch := make(chan SearchResultWithError)

	go func() {
		searchResult, err := sift.ExecuteSiftCmd(cmd)
		ch <- SearchResultWithError{searchResult, err}
	}()

	select {
	case resultWithError := <-ch:
		return resultWithError.SearchResult, resultWithError.err
	case <-time.After(time.Millisecond * time.Duration(timeout)):
		return sift.SearchResult{}, errors.New(fmt.Sprintf("sift search timeout for %dms limit", timeout))
	}
}

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


func testLinesRead() {
	lines, err := getFileLines("G:/raid/youxin.357.com/logs/main/log20180830", 8000, 8600)
	if err != nil {
		panic(err)
	}
	fmt.Println(lines)
	fmt.Println(len(lines))
	fmt.Println(lines[0])
}

func testSiftSearch() {
	searchResult, err := ExecuteSiftCmdWithTimeout("-e sift main.go -n --follow", 15 * 1000)

	fmt.Println(123)

	if err == nil {
		for _, result := range(searchResult.Results) {
			fmt.Println("这是一个文件的搜索结果：")
			sift.PrintResult(result)
			fmt.Println(result.Target)
			for _, match := range(result.Matches) {
				fmt.Println(match)
			}

			fmt.Println(result.Matches[0])
		}
		fmt.Println("执行耗时", searchResult.TimeCost)
	} else {
		fmt.Println("执行错误", err)
	}
}

type FileTreeNode struct {
	root string
	filePaths []string
	dirPaths []string
}

var (
	nodeChan chan *FileTreeNode
)

// TODO 关闭chan的时机

func getDirTree(dirPath string)  {

	nodeChan = make(chan *FileTreeNode, 1)

	go getFileTreeNode(dirPath)

	for node := range nodeChan {
		fmt.Println(node.root)
		for _, subDirPath := range(node.dirPaths) {
			go getFileTreeNode(subDirPath)
		}
	}

}

func getFileTreeNode(dirPath string) {
	filePaths, dirPaths, err := getDirContents(dirPath, 0)
	if err != nil {
		fmt.Println(err);
		os.Exit(1)
	}

	nodeChan <- &FileTreeNode{
		root: dirPath,
		filePaths: filePaths,
		dirPaths: dirPaths,
	}

	// close(nodeChan)
}

// 获取目录中的子目录和文件
func getDirContents(dirPath string, limit int64) ([]string, []string, error) {
	var filePaths []string
	var subDirPaths []string
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		//if info.IsDir() && info.Name() == subDirToSkip {
		//	fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
		//	return filepath.SkipDir
		//}
		if info.IsDir() && path != dirPath {
			subDirPaths = append(subDirPaths, path)
		} else {
			filePaths = append(filePaths, path)
		}
		return nil
	})
	return filePaths, subDirPaths, err
}