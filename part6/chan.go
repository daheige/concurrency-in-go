package main

import (
	"fmt"
	"runtime"
	"time"
)

var ch = make(chan struct{})

//多任务资源竞争
func printNum1(n int) {
	for i := 0; i < n; i++ {
		fmt.Println("ai=", i)
		time.Sleep(200 * time.Millisecond)
	}

	ch <- struct{}{}
}

func printNum2(n int) {
	for i := 0; i < n; i++ {
		fmt.Println("bi=", i)
		time.Sleep(200 * time.Millisecond)
	}

	<-ch //当管道有数据或者关闭了chan，才会释放
}

func main() {
	go printNum1(10)
	go printNum2(10)

	for {
		runtime.Gosched() //让出cpu执行权限给其他goroutine
	}
}
