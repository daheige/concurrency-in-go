package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	// 探讨slice,map并发写
	var s []int
	var m = make(map[int]bool, 10)

	// 推荐使用chan 把独立携程中执行的结果放入chan中，而不是通过共享变量的方式进行赋值
	res := make(chan int, 10)
	var wg sync.WaitGroup
	var lock sync.RWMutex // 通过读写锁保证map并发写没有data race

	wg.Add(20)
	for i := 0; i < 10; i++ {
		//i := i
		go func(i int) {
			defer wg.Done()

			log.Println("i = ", i)
			s = append(s, i) // 不推荐使用这种方式，这种方式是在多个goroutine之间对s进行操作，共享了s

			// 对于m赋值
			lock.Lock()
			log.Println("hello")
			m[i] = true
			lock.Unlock()

			// 推荐使用通道的方式，把数据放入res通道中就可以
			res <- i
		}(i)

		go func(i int) {
			defer wg.Done()

			log.Println("current i = ", i)

			lock.Lock()
			log.Println("hai")
			m[i] = true // 对于map存在并发写，产生了数据竞争
			// 当这里不加锁的话，使用go run -race t.go,发现这里存在go race
			lock.Unlock()
		}(i)
	}

	wg.Wait()

	// 关闭通道写入操作，也就是关闭发送者
	close(res)

	log.Println("s =", s)

	for v := range res {
		log.Println("current v = ", v)
	}

	log.Println("m = ", m)
	log.Println("exec end")
	log.Println(111)
	log.Println(123)
	time.Sleep(time.Second)
}

/**
当没有对m进行加锁保护map并发读写，容易出现data race
2020/04/29 23:30:26 i =  1
==================
WARNING: DATA RACE
Read at 0x00c00012a020 by goroutine 9:
  main.main.func1()
      /Users/heige/web/go/data-race.go:26 +0x10b

解决方式

方式1:使用互斥锁sync.Mutex
方式2:使用chan管道

使用管道的效率要比互斥锁高,也符合Go语言的设计思想
执行结果
2020/04/29 23:40:05 current i =  4
2020/04/29 23:40:05 hai
2020/04/29 23:40:05 i =  0
2020/04/29 23:40:05 hello
2020/04/29 23:40:05 current i =  9
2020/04/29 23:40:05 hai
2020/04/29 23:40:05 current i =  2
2020/04/29 23:40:05 i =  1
2020/04/29 23:40:05 current i =  1
2020/04/29 23:40:05 current i =  0
2020/04/29 23:40:05 i =  2
2020/04/29 23:40:05 i =  8
2020/04/29 23:40:05 i =  4
2020/04/29 23:40:05 hai
2020/04/29 23:40:05 hello
2020/04/29 23:40:05 hello
2020/04/29 23:40:05 hai
2020/04/29 23:40:05 current i =  7
2020/04/29 23:40:05 hai
2020/04/29 23:40:05 current i =  5
2020/04/29 23:40:05 hai
2020/04/29 23:40:05 i =  3
2020/04/29 23:40:05 hai
2020/04/29 23:40:05 current i =  3
2020/04/29 23:40:05 hello
2020/04/29 23:40:05 i =  5
2020/04/29 23:40:05 hai
2020/04/29 23:40:05 i =  6
2020/04/29 23:40:05 hello
2020/04/29 23:40:05 hello
2020/04/29 23:40:05 current i =  6
2020/04/29 23:40:05 hello
2020/04/29 23:40:05 hello
2020/04/29 23:40:05 current i =  8
2020/04/29 23:40:05 hai
2020/04/29 23:40:05 i =  9
2020/04/29 23:40:05 i =  7
2020/04/29 23:40:05 hai
2020/04/29 23:40:05 hello
2020/04/29 23:40:05 hello
2020/04/29 23:40:05 s = [0 1 2 8 4 3 5 6 9 7]
2020/04/29 23:40:05 current v =  0
2020/04/29 23:40:05 current v =  4
2020/04/29 23:40:05 current v =  1
2020/04/29 23:40:05 current v =  3
2020/04/29 23:40:05 current v =  2
2020/04/29 23:40:05 current v =  8
2020/04/29 23:40:05 current v =  5
2020/04/29 23:40:05 current v =  6
2020/04/29 23:40:05 current v =  9
2020/04/29 23:40:05 current v =  7
2020/04/29 23:40:05 m =  map[0:true 1:true 2:true 3:true 4:true 5:true 6:true 7:true 8:true 9:true]
2020/04/29 23:40:05 exec end
*/
