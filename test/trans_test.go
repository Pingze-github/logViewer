package main

import "fmt"

// 测试结构体的地址传递

type A struct {
	val int64
}

func transValue(a A) {
	a.val = 2
}

func transAddress(a *A) {
	(*a).val = 2
}

func transAddress2(aP *A) {
	// 错误，这样a还是复制的
	a := *aP
	a.val = 3
}

func main1() {
	a := A{val:1}
	transValue(a)
	fmt.Println(a)
	transAddress(&a)
	fmt.Println(a)
	transAddress2(&a)
	fmt.Println(a)
}
