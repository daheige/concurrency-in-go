//控制并发数
// 有时需要定时执行几百个任务，例如每天定时按城市来执行一些离线计算的任务。但是并发数又不能太高，
// 因为任务执行过程依赖第三方的一些资源，对请求的速率有限制。这时就可以通过 channel 来控制并发数。

package main

import (
	"log"
	"math/rand"
	"time"
)

// 构建一个缓冲型的 channel，容量为 4。接着遍历任务列表，每个任务启动一个 goroutine 去完成。
// 真正执行任务，访问第三方的动作在 w() 中完成，在执行 w() 之前，先要从 limit 中拿“许可证”，
// 拿到许可证之后，才能执行 w()，并且在执行完任务，要将“许可证”归还。
// 这样就可以控制同时运行的 goroutine 数。
// 这里，ch <- struct{}{} 放在 func 内部而不是外部，原因是：
// 如果在外层，就是控制系统 goroutine 的数量，可能会阻塞 for 循环，影响业务逻辑。
// 当缓冲区满了，就会进入等待状态，直到缓冲区有空闲就会执行w()
// ch 其实和逻辑无关，只是性能调优，放在内层和外层的语义不太一样。
// 还有一点要注意的是，如果 w() 发生 panic，那“许可证”可能就还不回去了，因此需要使用 defer 来保证

func main() {
	//定义一个缓冲区，限制goroutine执行频率
	//开辟空结构体，用来处理数据的发送处理

	t := time.Now()
	gNums := 1000
	total := 10 * 10000
	ch := make(chan struct{}, gNums) //缓冲区
	// done := make(chan bool, total)   //bool占据1字节，int占据4字节或8字节
	done := make(chan struct{}, total) //bool占据1字节，int占据4字节或8字节

	emptyStruct := struct{}{} //复用空结构体,这里基本上不占据内存

	for i := 0; i < total; i++ {
		// ch <- struct{}{}
		go func(i int) {
			defer func() {
				if err := recover(); err != nil {
					log.Println("exec panic")
					log.Println("panic info: ", err)
				}

				<-ch //必须在这里取出和放入done
				// done <- true
				done <- emptyStruct
			}()

			//该句建议放在goroutine里面
			// ch <- struct{}{}
			ch <- emptyStruct

			w()
			log.Println("current index: ", i)
			log.Println("hello")
		}(i)
	}

	//取出done chan
	for n := 0; n < total; n++ {
		log.Println("n = ", n)
		<-done
	}

	log.Println("cost time: ", time.Now().Sub(t).Seconds())

}

func w() {
	t := time.Now().UnixNano()
	n := time.Duration(rand.Int63n(1000))
	rand.Seed(t)
	log.Println("working...")
	//随机停顿
	time.Sleep(n * time.Millisecond)
	//模拟程序出现panic
	if n > 600 {
		panic("haha")
	}
}

/**
2019/07/25 21:35:59 n =  99999
2019/07/25 21:35:59 exec panic
2019/07/25 21:35:59 panic info:  haha
2019/07/25 21:35:59 true
2019/07/25 21:35:59 cost time:  50.67913431
*/
