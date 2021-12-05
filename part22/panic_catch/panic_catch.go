package main

import (
	"log"
	"time"
)
/**
panic捕获是要看处在什么样的上下文范围中
一个原则：携程运行要看处在什么样的上下文范围中，尽量在独立携程中捕获panic，因为panic不能跨携程捕获
*/
func main() {
	log.Println("start task...")

	// 这里可以捕获main主携程的panic,并不能捕获doWork()()触发的panic
	// defer catchRecover()
	doWork()()

	done := make(chan struct{}, 1)
	go func() {
		defer close(done)

		start := time.Now()
		var i int
		for {
			if time.Now().Unix()-start.Unix() >= 5*60 {
				log.Println("main will exit")
				break
			}

			doWork()() // 这里调用，看是否可以捕获函数运行的panic操作

			log.Println("current index: ", i)
			i++

			time.Sleep(100 * time.Millisecond)
		}
	}()
	<-done

	log.Println("main exit success...")
}

// Task task func.
type Task func()

func catchRecover() {
	if err := recover(); err != nil {
		log.Println("err: ", err)
	}
}

func doWork() Task {
	return func() {
		// 如果是放在函数中执行的话，这里需要进行对本函数出现的panic进行捕获
		// 所以panic捕获是要看处在什么样的上下文范围中
		// 一个原则，携程运行要看处在什么样的上下文范围中，尽量在独立携程中捕获panic，因为panic不能跨携程捕获
		defer catchRecover()

		log.Println(1234)

		// 通过struct{} chan信号量实现携程同步操作
		done := make(chan struct{}, 1)
		go func() {
			// 需要在当前开辟的携程中进行捕获panic，进行recover恢复
			// 当前携程执行的上下文在Task
			defer catchRecover()
			defer close(done)

			log.Println("goroutine task exec...")
			panic(112) // 模拟panic操作
		}()
		<-done // 当done中没有值或done没有关闭之前，这里会一直等，阻塞在这里
		panic(111)

		log.Println("task finish")
	}
}

/**
2020/11/04 22:25:53 123
2020/11/04 22:25:53 err:  111
2020/11/04 22:25:53 current index:  2915
2020/11/04 22:25:53 11
2020/11/04 22:25:53 123
2020/11/04 22:25:53 err:  112
2020/11/04 22:25:53 1223
2020/11/04 22:25:53 123
2020/11/04 22:25:53 err:  111
2020/11/04 22:25:53 current index:  2916
2020/11/04 22:25:54 main will exit
2020/11/04 22:25:54 main exit success...
*/
