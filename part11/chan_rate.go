package main

import (
	"fmt"
	"time"
)

/**
channel频率控制
在对channel进行读写的时，go还提供了非常人性化的操作
那就是对读写的频率控制，通过time.Ticke实现
 */

func main(){
	//生产者
	requests:= make(chan int ,5)
	for i:=0;i<5;i++{
		requests<-i
	}

	close(requests)

	//消费者限定时间
	//每隔1s消费一个任务
	limiter := time.Tick(time.Second*1)
	for req:=range requests{
		<-limiter
		fmt.Println("requets",req,time.Now()) //执行到这里，需要隔1秒才继续往下执行，time.Tick(timer)上面已定义
	}

}
