package main

import (
	"fmt"
	"runtime"
)

func test() {
	defer fmt.Println("test run end")
	runtime.Goexit()      //终止当前运行的所在携程
	fmt.Println("abcabc") //这里不会执行
}

func main() {
	fmt.Println("cpu核数: ", runtime.NumCPU())
	fmt.Println(runtime.GOMAXPROCS(2))
	go func() {
		fmt.Println(1)
		test()
		fmt.Println("hello") //hello不会执行
	}()

	//死循环让主main不退出
	for {
		runtime.Gosched()
	}
}
