package main

import (
	"net/http"

	"golang.org/x/time/rate"
)

//http全局限流
//算法描述：用户配置的平均发送速率为r，则每隔1/r秒一个令牌被加入到桶中（每秒会有r个令牌放入桶中），
// 桶中最多可以存放b个令牌。如果令牌到达时令牌桶已经满了，那么这个令牌会被丢弃
// NewLimiter returns a new Limiter that allows events up to rate r and permits
// bursts of at most b tokens.
// 每秒有r个令牌放入桶中，桶中最大放入b个令牌
// 其中r表示速率，b表示桶中最大的数量
var limiter = rate.NewLimiter(50, 600)

func limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if limiter.Allow() == false {
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

//访问地址： http://localhost:4000/
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", okHandler)
	// Wrap the servemux with the limit middleware.
	http.ListenAndServe(":4000", limit(mux))
}

func okHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

/*
虽然在某些情况下使用单个全局速率限制器非常有用，但另一种常见情况是基于IP地址或API密钥等标识符
为每个用户实施速率限制器。我们将使用IP地址作为标识符
*/
