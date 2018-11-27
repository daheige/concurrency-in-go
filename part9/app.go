package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("catch error: ", err)
		}
	}()

	// 模拟独立携程出现panic，是否导致main主携程退出
	var wg sync.WaitGroup
	wg.Add(1)
	ch := make(chan int)
	go func() {
		defer func() {
			// if err := recover(); err != nil {
			// 	log.Println("catch error: ", err)
			// }
			wg.Done()
			close(ch) //关闭ch
		}()
		select {
		case ch <- 1:
			log.Println("send success")
		default:
			rand.Seed(time.Now().Unix())
			rnd := rand.Intn(100)
			if rnd > 50 {
				panic("current goroutine panic")
			} else {
				log.Println("rnd is ", rnd)
			}
		}

		log.Println(111)
	}()

	wg.Wait()

	for v := range ch {
		log.Println("read ch is: ", v)
	}

	log.Println("waiting...")

}

/**
 * 多次运行后，抛出异常panic,会导致整个main退出
panic: current goroutine panic

goroutine 5 [running]:
main.main.func2(0xc000018140, 0xc00001a180)
	/mygo/src/projects/demo/app.go:36 +0x214
created by main.main
	/mygo/src/projects/demo/app.go:21 +0xb7
exit status 2
*/
