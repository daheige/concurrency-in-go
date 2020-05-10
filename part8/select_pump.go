package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go pump1(ch1)
	go pump2(ch2)

	go suck(ch1, ch2)

	time.Sleep(1 * time.Second)
}

func pump1(ch chan int) {
	for i := 0; ; i++ {
		ch <- i * 2
	}
}

func pump2(ch chan int) {
	for i := 0; ; i++ {
		ch <- i + 5
	}
}

func suck(ch1, ch2 chan int) {

	// select 随机选择一个chan读取
	for {
		select {
		case v := <-ch1:
			fmt.Printf("Received on channel 1: %d\n", v)
		case v := <-ch2:
			fmt.Printf("Received on channel 2: %d\n", v)
		}
	}
}

/**
Received on channel 1: 428158
Received on channel 2: 217435
Received on channel 2: 217436
Received on channel 2: 217437
Received on channel 2: 217438
Received on channel 2: 217439
Received on channel 2: 217440
Received on channel 1: 428160
Received on channel 2: 217441
Received on channel 2: 217442
Received on channel 2: 217443
Received on channel 1: 428162
*/
