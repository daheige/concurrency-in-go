# golang rate包实现令牌桶限流
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

# golang rate包 ip限流
    参考ip_rate.md
    
# 令牌桶实战解说

    https://github.com/didip/tollbooth

# 支持的框架

    http,gin,chi,echo,httpRouter

# tollbooth http 限流

    package main

    import (
        "log"
        "net/http"

        "github.com/didip/tollbooth"
    )

    func HelloHandler(w http.ResponseWriter, req *http.Request) {
        w.Write([]byte("Hello, World!"))
    }

    func main() {
        // Create a request limiter per handler.
        //每秒100个
        limiter := tollbooth.NewLimiter(100, nil)

        http.Handle("/", tollbooth.LimitFuncHandler(limiter, HelloHandler))

        address := ":4000"
        log.Println("server has run: ", address)
        http.ListenAndServe(address, nil)
    }


# gin 限流

    https://github.com/didip/tollbooth_gin
