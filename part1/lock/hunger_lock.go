package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup     //wg计数器
var sharedLock sync.Mutex //共享锁标志

const runtime = 1 * time.Second

func main() {
	//贪婪的goroutine
	greedyWorker := func() {
		defer wg.Done()

		var count int
		for begin := time.Now(); time.Since(begin) <= runtime; {
			sharedLock.Lock()
			time.Sleep(3 * time.Nanosecond) //一次性停顿3ns会占用别的goroutine的资源,会出现资源抢夺的情况
			sharedLock.Unlock()
			count++
		}

		fmt.Printf("Greedy worker was able to execute %v work loops\n", count)
	}
	//知足的goroutine按需要执行
	politeWorker := func() {
		defer wg.Done()

		var count int
		for begin := time.Now(); time.Since(begin) <= runtime; {
			//通过共享锁按需停顿
			sharedLock.Lock()
			time.Sleep(1 * time.Nanosecond)
			sharedLock.Unlock()

			sharedLock.Lock()
			time.Sleep(1 * time.Nanosecond)
			sharedLock.Unlock()

			sharedLock.Lock()
			time.Sleep(1 * time.Nanosecond)
			sharedLock.Unlock()

			count++
		}
		fmt.Printf("Polite worker was able to execute %v work loops.\n", count)
	}

	//开启2个wg计数器保证贪婪goroutine和知足goroutine执行完毕
	wg.Add(2)
	go greedyWorker()
	go politeWorker()
	wg.Wait()
}

/**
Polite worker was able to execute 221098 work loops.
Greedy worker was able to execute 284270 work loops
*/
