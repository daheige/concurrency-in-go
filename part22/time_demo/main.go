package main

import (
	"context"
	"log"
	"time"
)

func main() {
	// AsyncCall(func() {
	// 	for i := 0; i < 100; i++ {
	// 		log.Println("current index: ", i)
	// 		time.Sleep(200 * time.Millisecond)
	// 	}
	// }, 1*time.Second)
	exit := make(chan struct{}, 1)
	AsyncCallByCtx(func() {
		for i := 0; i < 100; i++ {
			select {
			case <-exit:
				log.Println("exit func call")
				return
			default:
			}

			log.Println("current index2: ", i)
			time.Sleep(300 * time.Millisecond)
		}
	}, 3*time.Second, exit)
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
func AsyncCallByCtx(fn func(), timeout time.Duration, exit chan struct{}) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "abc", 123)
	ctx, cancel := context.WithTimeout(ctx, timeout)
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

		close(exit) // 退出信号量,用来通知fn func函数中断执行，这里仅仅让fn func for逻辑不再执行了

		// 当上下文被取消后，Value是可以拿到ctx数据
		log.Println("abc: ", ctx.Value("abc")) // abc:  123 这里从上下文获取key
		log.Println("timeout")
	case <-done:
		log.Println("call success")
		return
	}
}
