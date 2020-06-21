// Package timeout
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// AsyncCall 第一种：使用time.NewTimer
func AsyncCall() {
	timer := time.NewTimer(1 * time.Second)
	done := make(chan struct{}, 1)
	go func() {
		defer close(done)

		// time.Sleep(300 * time.Millisecond)

		time.Sleep(1200 * time.Millisecond)

		log.Println("send msg: hello")
	}()

	select {
	case <-done:
		// 监听到Done信号，然后关闭定时器
		timer.Stop()
		log.Println("call success")
	case <-timer.C:
		log.Println("timeout")
	}
}

// AsyncCallWithCtx do task with ctx
func AsyncCallWithCtx(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	done := make(chan struct{}, 1)
	go func(ctx context.Context) {

		// time.Sleep(300 * time.Millisecond)

		time.Sleep(3000 * time.Millisecond)

		log.Println("send msg: hello")
		close(done)
	}(ctx)

	select {
	case <-done:
		log.Println("call success")
	case <-ctx.Done():
		log.Println("timeout")
		log.Println("timeout reason: ", ctx.Err())
	}
}

func main() {
	// AsyncCall()
	AsyncCallWithCtx(context.Background())

	// 如果把下面的信号量监控机制全部打开，对于http 服务来说，依然会执行send msg: hello,其实这里的超时限制，可以对于mysql,redis,mongodb 这样的服务调用设置超时
	// 比如调用http api接口，调用方可以设置超时，而不是服务端超时限制，服务端超时这个是一个伪命题，因为虽然提前返回了，但是后续的操作依然会执行
	ch := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// recivie signal to exit main goroutine
	// window signal
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)

	// linux signal,please use this in production.
	// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2, os.Interrupt, syscall.SIGHUP)

	// Block until we receive our signal.
	sig := <-ch

	log.Println("exit signal: ", sig.String())

	// 5s之后平滑退出
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	<-ctx.Done()

	log.Println("shutting down")
}

/**
% go run timeout.go
2020/06/21 17:40:12 timeout
2020/06/21 17:40:12 timeout reason:  context deadline exceeded
2020/06/21 17:40:13 send msg: hello
^C2020/06/21 17:40:46 exit signal:  interrupt
2020/06/21 17:40:51 shutting down
*/
