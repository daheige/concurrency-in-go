package main

/**
关于生产者和消费者平滑退出
对于一些job,web,grpc服务来说，需要控制平滑退出
下面的demo，演示了如何平滑退出程序
*/

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var wait = 2 * time.Second

func main() {
	//平滑重启
	ch := make(chan os.Signal, 1)
	exit := make(chan struct{}, 1)

	//将生产者和消费者的执行权限缩小，放在局部范围内执行，也就是缩小上下文的范围
	go proAndConsume(exit)

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// recivie signal to exit main goroutine
	// window signal
	// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2, os.Interrupt, syscall.SIGHUP)

	// Block until we receive our signal.
	sig := <-ch

	//关闭通道
	close(exit)

	log.Println("exit signal: ", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	<-ctx.Done()

	log.Println("server will exit...")
}

func proAndConsume(exit chan struct{}) {
	data := make(chan string, 100)

	//多个生产者
	for i := 0; i < 2; i++ {
		go func() {
			for {
				select {
				case <-exit: //多次读取关闭的通道，返回零值
					log.Println("production will exit...")
					return
				default:
					log.Println("work...")

					for i := 0; i < 2; i++ {
						data <- "hello" + strconv.Itoa(i)
					}
				}

			}
		}()
	}

	//多个消费者
	for i := 0; i < 2; i++ {
		go func() {
			for k := range data {
				log.Println("k: ", k)
			}
		}()
	}

}

/**
按下ctrl+c退出
2019/09/03 23:47:51 k:  hello0
2019/09/03 23:47:51 work...
2019/09/03 23:47:51 production will exit...
2019/09/03 23:47:51 k:  hello1
2019/09/03 23:47:51 work...
2019/09/03 23:47:51 production will exit...
2019/09/03 23:47:51 k:  hello0
2019/09/03 23:47:51 k:  hello1
2019/09/03 23:47:51 k:  hello0
2019/09/03 23:47:53 server will exit...
*/
