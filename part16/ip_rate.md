# ip 限流

    package main
    import (
        "net/http"
        "sync"
        "time"
        "golang.org/x/time/rate"
    )
    // Create a custom visitor struct which holds the rate limiter for each
    // visitor and the last time that the visitor was seen.
    type visitor struct {
        limiter  *rate.Limiter
        lastSeen time.Time
    }

    // Change the the map to hold values of the type visitor.
    var visitors = make(map[string]*visitor)
    var mtx sync.Mutex
    // Run a background goroutine to remove old entries from the visitors map.
    func init() {
        go cleanupVisitors()
    }

    func addVisitor(ip string) *rate.Limiter {
        limiter := rate.NewLimiter(2, 5)
        mtx.Lock()
        // Include the current time when creating a new visitor.
        visitors[ip] = &visitor{limiter, time.Now()}
        mtx.Unlock()
        return limiter
    }
    func getVisitor(ip string) *rate.Limiter {
        mtx.Lock()
        v, exists := visitors[ip]
        if !exists {
            mtx.Unlock()
            return addVisitor(ip)
        }
        // Update the last seen time for the visitor.
        v.lastSeen = time.Now()
        mtx.Unlock()
        return v.limiter
    }

    // Every minute check the map for visitors that haven't been seen for
    // more than 3 minutes and delete the entries.
    func cleanupVisitors() {
        for {
            time.Sleep(time.Minute)
            mtx.Lock()
            for ip, v := range visitors {
                if time.Now().Sub(v.lastSeen) > 3*time.Minute {
                    delete(visitors, ip)
                }
            }
            mtx.Unlock()
        }
    }

    func limit(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            limiter := getVisitor(r.RemoteAddr)
            if limiter.Allow() == false {
                http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
