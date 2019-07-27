package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello,world"))
	})

	addr := "0.0.0.0:8080"
	s := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Println("server has run on: ", addr)

	//s.ListenAndServe()启动，底层是一个Serve方法,里面是每个请求都是一个goroutine在处理
	// 内部是一个for循环监听，所以这里采用goroutine来启动，这样在main携程后面还可以捕捉退出信号
	// 当程序收到了退出信号后，就会调用s.Shutdown方法，优雅的退出http服务
	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Println("server error: ", err)
		}
	}()

	//优雅的关闭http服务
	//声明一个退出信号
	ch := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// recivie signal to exit main goroutine
	//window signal
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)
	// linux signal
	// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2, os.Interrupt, syscall.SIGHUP)

	// Block until we receive our signal.
	sig := <-ch

	log.Println("exit signal: ", sig.String())

	//5秒后自动退出
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//优雅的退出http服务
	e := make(chan error, 1) //接收退出的错误
	go func() {
		e <- s.Shutdown(ctx)
	}()

	<-ctx.Done()
	log.Println("shutdown error: ", <-e)
	log.Println("server will exit")
}
