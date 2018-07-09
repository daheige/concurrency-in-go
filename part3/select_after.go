package main

import (
	"fmt"
	"time"
)

//常见的超时处理
func main() {
	var c <-chan int
	select {
	case <-c: //读取c通道的值,没有读取就会超时处理
	case <-time.After(1 * time.Second):
		fmt.Println("Timed out.")
	}

	//default操作
	start := time.Now()
	var c1, c2 <-chan int
	select {
	case <-c1:
	case <-c2:
	default:
		fmt.Printf("In default after %v\n\n", time.Since(start))
	}
	//In default after 3.956µs
}
