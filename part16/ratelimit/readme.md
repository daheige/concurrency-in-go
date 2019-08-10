# 令牌桶实战

    https://github.com/didip/tollbooth

# 支持的框架

    http,gin,chi,echo,httpRouter

# go-chi/chi 限流

    package main

    import (
        "github.com/didip/tollbooth"
        "github.com/didip/tollbooth_chi"
        "github.com/pressly/chi"
        "net/http"
        "time"
    )

    func main() {
        // Create a limiter struct.
        limiter := tollbooth.NewLimiter(1, time.Second, nil)

        r := chi.NewRouter()

        r.Use(tollbooth_chi.LimitHandler(limiter))

        r.Get("/", func(w http.ResponseWriter, r *http.Request) {
            w.Write([]byte("Hello, world!"))
        })

        http.ListenAndServe(":12345", r)
    }

# gin 限流

    https://github.com/didip/tollbooth_gin
