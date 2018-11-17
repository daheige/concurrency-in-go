package main

import (
	"log"
	"runtime"
	"sync"
	"thinkgo/common"
)

var count = 1

func main() {
	log.Println("fefe")

	var wg sync.WaitGroup

	//抢占式的更新count，需要对count进行枷锁保护
	//如果不加锁，count每次执行后，值都不一样
	lock := common.NewChanLock()
	for i := 0; i < 1000; i++ {
		runtime.Gosched() //让出当前cpu给其他goroutine执行

		wg.Add(1)
		go func() {
			defer wg.Done()
			lock.Lock()
			defer lock.Unlock()
			v := count
			log.Println("current count: ", v)
			v++
			count = v
		}()
	}
	log.Println("exec running....")
	wg.Wait()

	log.Println("count: ", count)
}
