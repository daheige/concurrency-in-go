package main

import (
	"fmt"
	"sync"
)

//利用wg保证goroutine执行完毕
func main() {
	var data int

	var wg sync.WaitGroup
	wg.Add(1) //通过wg的信号计数器来保证goroutine执行
	go func() {
		defer wg.Done()
		data++

	}()
	wg.Wait() //这里会一直等待goroutine执行完毕
	if data == 1 {
		fmt.Printf("the value is %v.\n", data)
	}

	fmt.Println(111)
}
