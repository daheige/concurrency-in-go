package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println("cpu核数: ", runtime.NumCPU()) //4
	fmt.Println(runtime.GOMAXPROCS(4))

	//死循环让主main不退出
	for {
		runtime.Gosched() //让出当前goroutine执行权限
		go func() {
			fmt.Println(1)
		}()

		fmt.Println(0)
	}
}
