/**
goroutine之间的通讯
goroutine本质上是协程，可以理解为不受内核调度，而受go调度器管理的线程
goroutine之间可以通过channel进行通信或者说是数据共享
当然你也可以使用全局变量来进行数据共享。
 */
package main

import (
	"fmt"
	"sync"
)

//生产者
func Productor(mychan chan int,data int,wait *sync.WaitGroup)  {
	mychan <- data
	fmt.Println("product data：",data)
	wait.Done()
}

//消费者
func Consumer(mychan chan int,wait *sync.WaitGroup)  {
	a := <- mychan
	fmt.Println("consumer data：",a)
	wait.Done()
}

func main() {
	datachan := make(chan int, 10)   //通讯数据管道
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go Productor(datachan, i,&wg) //生产数据
	}

	for j := 0; j < 10; j++ {
		wg.Add(1)
		go Consumer(datachan,&wg)  //消费数据
	}

	wg.Wait()
}

