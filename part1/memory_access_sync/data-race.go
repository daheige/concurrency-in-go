package main

import (
	"log"
	"runtime"
)

/**
关于竞争状态：
如果两个或者多个 goroutine 在没有互相同步的情况下,访问某个共享的资源,并试图同时
读和写这个资源,就处于相互竞争的状态,这种情况被称作竞争状态(race candition)。竞争状态
的存在是让并发程序变得复杂的地方,十分容易引起潜在问题。对一个共享资源的读和写操作必
须是原子化的,换句话说,同一时刻只能有一个 goroutine 对共享资源进行读和写操作。

每个 goroutine 都会覆盖另一个 goroutine 的工作。这种覆盖发生在 goroutine 切换的时候。每
个 goroutine 创造了一个 counter 变量的副本,之后就切换到另一个 goroutine。当这个 goroutine
再次运行的时候, counter 变量的值已经改变了,但是 goroutine 并没有更新自己的那个副本的
值,而是继续使用这个副本的值,用这个值递增,并存回 counter 变量,结果覆盖了另一个
goroutine 完成的工作。
*/

var counter = 0

func main() {
	done := make(chan struct{}, 4)
	for i := 0; i < 4; i++ {
		go func(i int) {
			log.Println("current index: ", i)
			val := counter
			runtime.Gosched()

			val++
			counter = val

			done <- struct{}{}

		}(i)
	}

	for i := 0; i < 4; i++ {
		<-done
	}

	//输出的结果可能是1 2 3 4
	log.Println("counter=", counter)
}
