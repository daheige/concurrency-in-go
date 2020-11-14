package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// http server服务中，对于每个处理器函数或处理器，一般建议对panic进行捕获实现recover恢复
// 因为panic不能跨携程进行捕获
// 如果因为一些没有捕获的panic在runtime过程中，会导致服务进程崩溃退出，也就是整个main进行退出
// 对于一些web框架，比如gin,gorilla/mux,go-chi/chi都是采用中间件的方式，在middleware中进行捕获
// 它们把所有的逻辑，进行二次封装，也就是说逻辑都是运行在handler中，而中间件中对runtime中的panic做了recover()恢复处理
// 具体gin demo: https://github.com/daheige/goapp/blob/master/internal/web/middleware/log.go#L62
// https://github.com/daheige/goapp/blob/master/internal/web/routes/web.go#L23
func main() {
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/", Index)
	httpMux.HandleFunc("/test", MockPanic)

	addr := fmt.Sprintf("0.0.0.0:%d", 1338)
	log.Println("http has run: ", addr)
	server := &http.Server{
		Handler: httpMux,
		Addr:    addr,
		// Good practice to set timeouts to avoid Slowloris attacks.
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 2 << 20, // header max 2MB
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln("err: ", err)
	}
}

// Index index.
func Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

// MockPanic test panic.
func MockPanic(w http.ResponseWriter, r *http.Request) {
	// 捕获异常的panic
	// 这里不会捕获独立携程中的panic.
	// 也就是无法捕获到http://localhost:1338/test?id=2 抛出的panic.
	// 这个defer 只能捕获id=1时候发生的panic,不能跨携程捕获id=2的请求
	// 对于下面的id=1请求才可以捕获到panic
	// http://localhost:1338/test?id=1
	// 2020/11/14 08:29:11 exec panic:  param error
	defer func() {
		if err := recover(); err != nil {
			log.Println("exec panic: ", err)
			w.Write([]byte("server inner error"))
		}
	}()

	id := r.FormValue("id")
	if id == "1" {
		// 程序抛出panic,一般建议对panic进行捕获处理
		panic("param error")
		// panic(http.ErrAbortHandler)
	}

	// 如果不对下面的独立携程进行panic捕获，整个main进程都退出了
	// http://localhost:1338/test?id=2 请求地址
	// 当服务端终端打印出如下信息后
	// 导致http server 进行退出，也就是进程退出
	// 2020/11/13 23:40:10 goroutine run...
	// panic: current goroutine panic!
	// 可是当再次访问http://localhost:1338/ ，整个main服务都崩溃了
	go func() {
		// 这里必须在独立携程的上下文中捕获panic，才不会导致main退出
		// 捕获这个独立携程中的panic
		/*defer func() {
			if err := recover(); err != nil {
				log.Println("goroutine exec panic: ", err)
			}
		}()*/

		log.Println("goroutine run...")
		panic("current goroutine panic!")
	}()

	w.Write([]byte("id=" + id))
}

/**
2020/11/13 23:16:28 http: panic serving [::1]:57402: param error
goroutine 6 [running]:
net/http.(*conn).serve.func1(0xc00004eaa0)
	d:/Go/src/net/http/server.go:1772 +0x140
panic(0x680500, 0x743180)
	d:/Go/src/runtime/panic.go:975 +0x3f1
main.MockPanic(0x74d780, 0xc00012a0e0, 0xc000134000)
	D:/web/go/demo/server/main.go:42 +0xfd
net/http.HandlerFunc.ServeHTTP(0x6fd718, 0x74d780, 0xc00012a0e0, 0xc000134000)
	d:/Go/src/net/http/server.go:2012 +0x4b
net/http.(*ServeMux).ServeHTTP(0xc0000382c0, 0x74d780, 0xc00012a0e0, 0xc000134000)
	d:/Go/src/net/http/server.go:2387 +0x1ac
net/http.serverHandler.ServeHTTP(0xc00012a000, 0x74d780, 0xc00012a0e0, 0xc000134000)
	d:/Go/src/net/http/server.go:2807 +0xaa
net/http.(*conn).serve(0xc00004eaa0, 0x74dc80, 0xc000038340)
	d:/Go/src/net/http/server.go:1895 +0x873
created by net/http.(*Server).Serve
	d:/Go/src/net/http/server.go:2933 +0x363
*/
