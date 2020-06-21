package main

/**
 * 采用context+select实现超时调用
 * 一般用在请求远端服务或接口，以及定时task的执行
 */
import (
	"context"
	"log"
	"time"
)

func stuffHandler() {
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// 跟上面的等价
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer cancel()

	done := make(chan struct{}, 1) //当任务完成后，接收done信号
	go doStuff(done)

	select {
	case <-ctx.Done():
		log.Printf("handler will exit...")
		return
	case <-done:
		log.Println("stuff task has done")
	}
}

func doStuff(done chan struct{}) {
	defer close(done)

	i := 0
	for {
		if i >= 1000 {
			return
		}

		time.Sleep(100 * time.Microsecond) //每次停顿100ms
		log.Println("hello: ", i)
		i++
	}
}

func main() {
	log.Println("task exec begin")
	// stuffHandler()
	handler(5)
	log.Println("task exec end")
	// for {
	// }
}

/**
把上面的for{}打开
% go run ctx.go
2020/06/21 17:31:11 task exec begin
2020/06/21 17:31:16 timeout
2020/06/21 17:31:16 task exec end
2020/06/21 17:31:21 1
2020/06/21 17:31:21 daheige
2020/06/21 17:31:21 hahaha
*/

//精准rpc调用
//启动一个协程并执行 RPC 调用，同时初始化一个超时定时器。然后在主协程中同时监听 RPC 完成事件信号以及定时器信号。
//如果 RPC 完成事件先到达，则表示本次 RPC 成功，否则，当定时器事件发生，表明本次 RPC 调用超时。
//这种模型确保了无论何种情况下，一次 RPC 都不会超过预定义的时间，实现精准控制超时
//具体业务场景，可以基于这个handler来封装自己的业务
func handler(t int) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t)*time.Second)
	defer cancel()

	//设置值到ctx上
	ctx = context.WithValue(ctx, "name", "daheige")

	//用chan保证goroutine执行完毕
	done := make(chan struct{}, 1)
	go func() {
		rpc(ctx, 1)

		log.Println("hahaha")
		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		log.Println("timeout")
	case <-done:
		log.Println("task has done")
	}
}

//模拟rpc调用
func rpc(ctx context.Context, i int) {
	time.Sleep(10 * time.Second) //停顿10s 打开模拟超时
	log.Println(i)
	log.Println(ctx.Value("name"))
}
