/**
由于goroutine是异步执行的，那很有可能出现主程序退出时还有goroutine没有执行完，
此时goroutine也会跟着退出。此时如果想等到所有goroutine任务执行完毕才退出，
go提供了sync包和channel来解决同步问题，当然如果你能预测每个goroutine执行的时间，
你还可以通过time.Sleep方式等待所有的groutine执行完成以后在退出程序(如上面的列子)。
 */
package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

/*
1. 使用sync包同步goroutine
	sync大致实现方式WaitGroup 等待一组goroutinue执行完毕. 主程序调用 Add 添加等待的goroutinue数量
	每个goroutinue在执行结束时调用 Done
	此时等待队列数量减1.，主程序通过Wait阻塞，直到等待队列为0.
2.采用chan实现同步
	通过channel能在多个groutine之间通讯，当一个goroutine完成时候向channel发送退出信号
	等所有goroutine退出时候，利用for循环channe去读取channel中的信号，若取不到数据会阻塞原理，
	等待所有goroutine执行完毕，使用该方法有个前提是你已经知道了你启动了多少个goroutine。
 */

func main(){
	var wg sync.WaitGroup //声明一个WaitGroup变量
	for i :=0 ; i<10 ;i++{
		wg.Add(1) // WaitGroup的计数加1
		go cal2(i,i+1,&wg)
	}

	wg.Wait()  //等待所有goroutine执行完毕
	log.Println("exec end")

	done := make(chan bool,10)  //声明并分配管道内存
	for i :=0 ; i<10 ;i++{
		go cal3(i,i+1,done)
	}

	for j :=0; j<10; j++{
		<- done  //取信号数据，如果取不到则会阻塞
	}

	close(done) // 关闭管道

}

func cal2(a int , b int ,wg *sync.WaitGroup)  {
	defer wg.Done() //goroutinue完成后, WaitGroup的计数-1

	c := a+b
	fmt.Printf("%d + %d = %d\n",a,b,c)
}

func cal3(a int , b int ,done chan bool)  {
	c := a+b
	fmt.Printf("cal3: %d + %d = %d\n",a,b,c)
	time.Sleep(time.Second*2)
	done <- true
}
