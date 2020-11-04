package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	done := make(chan struct{}, 1)
	go func() {
		defer close(done)

		log.Println(111)
	}()

	<-done

	// time.Sleep(1 * time.Second)
	// go version >=1.14.x版本后，这里会发生携程阻塞，死锁，从而panic
	// 这种做法也是不可取
	// select {}

	// 下面的做法也不可取，一般采用接收中断信号量方式，防止main主携程退出
	/*ch := make(chan struct{}, 1)
	for {
		select {
		case <-ch:
		}
	}*/

	ch := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// recivie signal to exit main goroutine
	// window signal
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)

	// linux signal if you use linux on production,please use this code.
	// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2, os.Interrupt, syscall.SIGHUP)

	// Block until we receive our signal.
	sig := <-ch

	log.Println("exit signal: ", sig.String())
	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Optionally
	<-ctx.Done()
	log.Println("main exit success")
}

/**
% go run app.go
2020/11/04 22:59:42 111
fatal error: all goroutines are asleep - deadlock!

go version >= 14.x.x版本后会发生panic，携程阻塞
goroutine 1 [select (no cases)]:
main.main()
        /Users/heige/web/go/demo/cmd/app.go:19 +0x85

// 采用中断信号量方式退出main
% go run app.go
2020/11/04 23:14:18 111
^C2020/11/04 23:14:22 exit signal:  interrupt
heige@daheige cmd % go run app.go
2020/11/04 23:14:58 111
^C2020/11/04 23:15:11 exit signal:  interrupt
2020/11/04 23:15:16 main exit success
*/
