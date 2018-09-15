package main

import (
	"errors"
	"fmt"
	"time"
	"github.com/Pingze-github/sift"
)

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


func main() {
	searchResult, err := ExecuteSiftCmdWithTimeout("-e sift . -n", 15 * 1000)

	fmt.Println(123)

	if err == nil {
		for _, result := range(searchResult.Results) {
			fmt.Println("这是一个文件的搜索结果：")
			sift.PrintResult(result)
		}
		fmt.Println("执行耗时", searchResult.TimeCost)
	} else {
		fmt.Println("执行错误", err)
	}

}
