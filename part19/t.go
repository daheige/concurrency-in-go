package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	log.Println("hello")

	// 一主多子协程的方式，推荐使用sync.WaitGroup
	// 单个通道进行通信，推荐使用chan通道方式实现，参考app.go
	var wg sync.WaitGroup
	var t = time.Now()
	wg.Add(3)
	go worker(&wg, 1000)
	go worker(&wg, 1000)
	go worker(&wg, 1000)
	wg.Wait()

	log.Println("cost time: ", time.Now().Sub(t))
}

func worker(wg *sync.WaitGroup, n int) {
	t := time.Now()
	sum := 0
	for i := 0; i < n; i++ {
		log.Println("i: ", i)
		sum += i
	}

	log.Println("worker time: ", time.Now().Sub(t))
	log.Println("sum = ", sum)

	wg.Done()
}
