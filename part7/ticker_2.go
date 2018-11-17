package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	//执行一次定时器操作
	tick := time.NewTicker(1 * time.Second)
	t := <-tick.C
	log.Println(111)
	log.Println("t: ", t)

	//每隔多久执行动作
	//time.Tick底层每次都new NewTicker一个对象
	for range time.Tick(time.Millisecond * 300) {
		fmt.Println("111")
	}
}
