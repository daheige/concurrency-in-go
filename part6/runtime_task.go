package main

import (
	"fmt"
	"runtime"
	"time"
)

//多任务资源竞争
func printNum(n int, g string) {
	for i := 0; i < n; i++ {
		fmt.Println(g, "i=", i)
		time.Sleep(200 * time.Millisecond)
	}
}
func main() {
	go printNum(100, "a")
	go printNum(100, "b")

	//主携程main不退出
	for {
		runtime.Gosched()
	}

}
