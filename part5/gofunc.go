package main

import "log"

func main() {
	ch := make(chan int)
	//生产者
	go test(ch)

	//消费者
	for c := range ch {
		log.Println("current ch: ", c)
	}

	log.Println(111)
}

func test(ch chan int) {
	for i := 0; i < 10; i++ {
		ch <- i
	}

	close(ch) //这里必须关闭，range读取的时候，才会释放
}
