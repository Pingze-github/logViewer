package main

import (
	"fmt"
	"sync"
	"time"
)

func test100() {

	wg := sync.WaitGroup{}

	stopChan := make(chan int)
	valsChan := make(chan int)

	go func () {
		// 1s后停止
		time.Sleep(time.Duration(10e9))
		stopChan <- 1
	}()

	go func() {
		for i:=0; i<3; i++ {
			wg.Add(1)
			valsChan <- i
		}

		wg.Wait()
		fmt.Println("wait finish")
		stopChan <- 1
	}()

	//for {
	//	select {
	//		case <- stopChan:
	//			fmt.Println("<- stopChan")
	//			return
	//		// case val := <- valsChan:
	//		default:
	//			fmt.Println("取出")
	//			// 最后一次取出后，循环卡在这里，等待valsChan输入值，形成死锁
	//			val := <- valsChan
	//			time.Sleep(time.Duration(1e9))
	//			fmt.Println("val", val)
	//			fmt.Println("b", wg)
	//			wg.Done()
	//			// 在这等待一下，可以保证wg.Wait()执行，使程序退出
	//			time.Sleep(time.Duration(1))
	//			fmt.Println(wg)
	//	}
	//}

	for val := range(valsChan) {
		fmt.Println("val", val)
	}

}

func test1 () {
	valsChan := make(chan int)

	go func() {
		for i := 0; i < 3; i++ {
			valsChan <- i
		}
	}()

	for val := range(valsChan) {
		fmt.Println("val", val)
	}
}

func test2 () {
	valsChan := make(chan int)
	stopChan := make(chan int)
	wg := sync.WaitGroup{}
	go func() {
		for i := 0; i < 3; i++ {
			wg.Add(1)
			valsChan <- i
		}
		wg.Wait()
		stopChan <- 1
	}()
	for {
		select {
		case <- stopChan:
			return
		case val := <- valsChan:
			wg.Done()
			fmt.Println("val", val)
		}
	}
}

func test3 () {
	valsChan := make(chan int)
	stopChan := make(chan int)
	wg := sync.WaitGroup{}
	go func() {
		for i := 0; i < 3; i++ {
			wg.Add(1)
			valsChan <- i
		}
		wg.Wait()
		stopChan <- 1
	}()
	for {
		select {
		case <- stopChan:
			return
		case val := <- valsChan:
			wg.Done()
			fmt.Println("val", val)
		}
	}
}

func main() {
	test3()
}

// select的所有case优先级相同，default比所有case优先级低
// FIXME default是非阻塞的，每次循环过来直接执行，不像case有阻塞条件