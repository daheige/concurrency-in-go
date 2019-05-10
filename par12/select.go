package main

import (
	"log"
	"math/rand"
	"time"
)

func main() {
	var c1, c2 = generator(), generator()
	//select随机选择
	tm := time.After(3 * time.Second) //运行3s后就退出程序
	for {
		//超时就退出
		select {
		case <-tm:
			log.Println("timeout")
			return
		default:
		}

		select {
		case n := <-c1:
			log.Println("received c1: ", n)
		case n := <-c2:
			log.Println("received c2: ", n)
		}
	}
}

//不停地产生int
func generator() chan int {
	out := make(chan int)
	go func() {
		i := 0
		for {
			time.Sleep(time.Duration(rand.Intn(120)) * time.Microsecond)
			out <- i
			i++
		}
	}()

	return out
}
