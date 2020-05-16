package main

import (
	"fmt"
	"time"
)

// 通过指定goroutine个数，work pool消费，防止goroutine泄露和暴涨
// 模拟生产者和消费者模式
// 下面是一个简易的work pool 复杂的work pool可以自行改变
func main() {
	jobsNum := 100
	jobs := make(chan int, jobsNum)
	results := make(chan int, jobsNum)
	// 开启多个个goroutine开始work pool消费
	for i := 0; i < 3; i++ {
		go worker(i, jobs, results)
	}

	// 把任务放入待消费的jobs队列中
	for j := 0; j < jobsNum; j++ {
		jobs <- j
	}

	// 关闭jobs通道
	close(jobs)

	// 打印输出结果
	for i := 0; i < jobsNum; i++ {
		fmt.Println("res = ", <-results)
	}
}

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("current worker index: ", id)
		fmt.Printf("worker id: %d start job:%d\n", id, j)
		// mock do something
		time.Sleep(200 * time.Millisecond)
		fmt.Printf("worker:%d end job:%d\n", id, j)

		results <- j * 2
	}
}

/**
worker id: 0 start job:93
worker:1 end job:95
worker:2 end job:94
current worker index:  2
worker id: 2 start job:97
res =  190
res =  188
worker:0 end job:93
current worker index:  0
worker id: 0 start job:98
current worker index:  1
worker id: 1 start job:96
res =  186
worker:1 end job:96
current worker index:  1
worker:0 end job:98
worker id: 1 start job:99
res =  192
res =  196
worker:2 end job:97
res =  194
worker:1 end job:99
res =  198
*/
