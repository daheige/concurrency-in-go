package main

import (
	"bytes"
	"fmt"
	"sync"
)

//非安全性并发操作
//在这个例子中，你可以看到，我们不需要通过通信同步内存访问或共享数据。
func main() {
	printData := func(wg *sync.WaitGroup, data []byte) {
		defer wg.Done()

		var buff bytes.Buffer
		for _, b := range data {
			fmt.Fprintf(&buff, "%c", b)
		}
		fmt.Println(buff.String())
	}

	//通过wg信号量计数器保证goroutine执行完毕
	var wg sync.WaitGroup
	wg.Add(2)
	data := []byte("golang")
	go printData(&wg, data[:3]) // 1
	go printData(&wg, data[3:]) // 2

	wg.Wait()
}
