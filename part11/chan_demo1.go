package main

import "fmt"

/**
channel俗称管道，用于数据传递或数据共享
其本质是一个先进先出的队列，使用goroutine+channel进行数据通讯简单高效
同时也线程安全，多个goroutine可同时修改一个channel，不需要加锁。

channel可分为三种类型：
只读channel：只能读channel里面数据，不可写入

只写channel：只能写数据，不可读

一般channel：可读可写

	var readOnlyChan <-chan int            // 只读chan
	var writeOnlyChan chan<- int           // 只写chan
	var mychan  chan int                     //读写channel
	//定义完成以后需要make来分配内存空间，不然使用会deadlock
	mychannel = make(chan int,10)

	//或者
	read_only := make (<-chan int,10)//定义只读的channel
	write_only := make (chan<- int,10)//定义只写的channel
	read_write := make (chan int,10)//可同时读写

读写数据需要注意的是：
	管道如果未关闭，在读取超时会则会引发deadlock异常
	管道如果关闭进行写入数据会pannic
	当管道中没有数据时候再行读取或读取到默认值，如int类型默认值是0
	ch <- "wd"  //写数据
	a := <- ch //读取数据
	a, ok := <-ch  //优雅的读取数据

循环管道需要注意的是：
	使用range循环管道，如果管道未关闭会引发deadlock错误。
	如果采用for死循环已经关闭的管道，当管道没有数据时候，读取的数据会是管道的默认值，并且循环不会退出。

带缓冲区channe和不带缓冲区channel
	带缓冲区channel：定义声明时候制定了缓冲区大小(长度)，可以保存多个数据。
	不带缓冲区channel：只能存一个数据，并且只有当该数据被取出时候才能存下一个数据
	ch := make(chan int) //不带缓冲区
	ch := make(chan int ,10) //带缓冲区。
*/
func main() {
	mychannel := make(chan int,10)
	for i := 0;i < 10;i++{
		mychannel <- i
	}

	close(mychannel)  //关闭管道
	fmt.Println("data lenght: ",len(mychannel))

	//从通道中读取数据
	for  v := range mychannel {  //循环管道
		fmt.Println(v)
	}

	fmt.Printf("data lenght:  %d",len(mychannel))
}


