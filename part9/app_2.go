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
		//在独立的goroutine中,如果出现了panic恐慌，需要对其捕获处理，否则将导致整个main退出
		defer func() {
			if err := recover(); err != nil {
				log.Println("catch error: ", err)
			}
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
