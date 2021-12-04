package main

import (
	"context"
	"log"
	"time"
)

func main() {
	AsyncCall(func() {
		for i := 0; i < 100; i++ {
			log.Println("current index: ", i)
			time.Sleep(200 * time.Millisecond)
		}
	}, 1*time.Second)

	AsyncCallByCtx(func() {
		for i := 0; i < 100; i++ {
			log.Println("current index2: ", i)
			time.Sleep(300 * time.Millisecond)
		}
	}, 3*time.Second)
	time.Sleep(2 * time.Second)
}

// AsyncCall fn func() 通过timer模式进行超时控制
func AsyncCall(fn func(), timeout time.Duration) {
	timer := time.NewTimer(timeout)
	done := make(chan struct{}, 1)
	go func() {
		defer close(done)

		fn()
	}()

	select {
	case <-done:
		log.Println("call success")
		return
	case <-timer.C:
		timer.Stop()
		log.Println("timeout")
		return
	}
}

// AsyncCallByCtx 通过ctx timeout方式实现fn调度，超时就取消fn执行
func AsyncCallByCtx(fn func(), timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	done := make(chan struct{}, 1)
	go func() {
		defer close(done)

		fn()
	}()

	select {
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded { // 表示ctx超时取消
			cancel()
		}

		log.Println("timeout")
	case <-done:
		log.Println("call success")
		return
	}
}
