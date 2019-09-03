package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var wait = 2 * time.Second

func main() {
	ch := make(chan os.Signal, 1)
	exit := make(chan struct{}, 1)

	//即使在异步生产中把数据放入data中，然后接收信号量退出，平滑退出
	//但这样无法做到平滑退出，这是一个误区
	//可以把下面的代码，放入一个局部范围内执行，缩小goroutine执行权限，拥有的上下文环境
	// 可以看app_graceful2.go
	data := make(chan string, 10)
	for i := 0; i < 3; i++ {
		go production(data, exit)
	}

	consume(data)

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// recivie signal to exit main goroutine
	// window signal
	// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2, os.Interrupt, syscall.SIGHUP)

	// Block until we receive our signal.
	sig := <-ch

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	<-ctx.Done()

	close(exit)

	log.Println("exit signal: ", sig.String())
	log.Println("server will exit...")
}

func production(data chan string, exit chan struct{}) {
	for {
		select {
		case <-exit:
			log.Println("production will exit...")
			return
		default:
			log.Println("work...")
		}

		data <- "hello"
	}
}

func consume(data chan string) {
	for d := range data {
		log.Println("recv data: ", d)
	}
}

//执行结果
/* 2019/09/03 23:19:08 recv data:  hello
2019/09/03 23:19:08 recv data:  hello
2019/09/03 23:19:08 recv data:  hello
2019/09/03 23:19:08 recv data:  hello
2019/09/03 23:19:08 work...
2019/09/03 23:19:08 work...
2019/09/03 23:19:08 work...
2019/09/03 23:19:08 work...
2019/09/03 23:19:08 work...
2019/09/03 23:19:08 work...
2019/09/03 23:19:08 work...
^Csignal: interrupt
并非平滑退出程序
*/
