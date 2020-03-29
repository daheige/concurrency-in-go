package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	// log使用方式
	log.Println("fefe")
	log.SetPrefix("[info]")
	log.Println(1111)
	log.SetFlags(log.LstdFlags | log.Llongfile)
	log.Println(12345)

	// 方式1.通过通道实现同步，在不同的goroutine之间传递信号，通过chan来进行通信
	// 这里的ch是一个非阻塞的缓冲通道
	ch := make(chan struct{}, 1)
	go func() {
		log.Println(123)
		time.Sleep(2 * time.Second)
		ch <- struct{}{}
	}()

	<-ch

	// 方式2： 通过sync.WaitGroup 协程计数器来进行同步
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done() // goroutine执行完毕后，计数器减1

		for i := 0; i < 2; i++ {
			log.Println("hello,index: ", i)
		}
	}()

	// 当计数为0，这里就解除了阻塞
	// Counter is 0, no need to wait.
	wg.Wait()

	log.Println(time.Now().Sub(start).Nanoseconds() / 1e3)
	log.Println(time.Now().Sub(start).Microseconds())

	log.Println(int64(time.Now().Sub(start) / time.Millisecond))
	log.Println("main will exit")
}
