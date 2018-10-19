package main

import (
	"fmt"
	"time"
)

func main() {
	//create a timer
	//只会响应一次
	timer := time.NewTicker(2 * time.Second)
	fmt.Println("当前时间", time.Now())

	//after 2s,input data to timer.C
	//util data,you can read data
	t := <-timer.C
	fmt.Println("t = ", t)

	t2 := time.NewTicker(3 * time.Second)
	go func() {
		<-t2.C
		//可以打印数据了
		fmt.Println(111)
	}()

	t2.Stop() //停止定时器

	//定时器重置
	t3 := time.NewTimer(3 * time.Second)
	t3.Reset(1 * time.Second)

	<-t3.C
	fmt.Println("时间到了")
	for {

	}
}
