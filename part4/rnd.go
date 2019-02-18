// 约定: 谁创建goroutine，谁负责停止
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	//产生随机数
	newRnd := func(done <-chan struct{}) <-chan int {
		rndCh := make(chan int)
		go func() {
			defer fmt.Println("new rnd int closure exited")

			for {
				select {
				case rndCh <- rand.Int():
					fmt.Println("rnd int has worked")
				case <-done: //接收到done信号(或者done通道北关闭）后，就退出当前goroutine
					return
				}
			}
		}()

		return rndCh
	}

	done := make(chan struct{})
	rndStream := newRnd(done) //在主main中创建
	fmt.Println("3 random its: ")
	for i := 0; i < 3; i++ {
		fmt.Printf("%d:%d\n", i, <-rndStream)
	}

	close(done) //关闭done信号
	fmt.Println("main will exit...")
	time.Sleep(1 * time.Second) //让主main不退出，保证goroutine运行

}
