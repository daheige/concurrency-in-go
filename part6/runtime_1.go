package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	var wg sync.WaitGroup //信号计数器，保证携程执行完毕
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			fmt.Println(i)
		}
	}()
	fmt.Println("等待执行完毕")
	wg.Wait()

	for i := 0; i < 100; i++ {
		runtime.Gosched() //让出cpu时间片，让出当前goroutine执行权
		// 调度器p会安排其他等待的任务执行，并在下一次某个时刻执行
		fmt.Println("主main的i: ", i)
	}
}
