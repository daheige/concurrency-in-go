package main

// select-case实现非阻塞channel
// 原理通过select+case加入一组管道，
// 当满足（这里说的满足意思是有数据可读或者可写)
// select中的某个case时候，那么该case返回，若都不满足case，则走default分支。

import (
	"fmt"
)

func send(c chan int)  {
	for i :=1 ; i<10 ;i++  {
		c <-i
		fmt.Println("send data : ",i)
	}
}

func main() {
	resch := make(chan int,20)
	strch := make(chan string,10)
	go send(resch)
	strch <- "wd"
	select {
	case a := <-resch:
		fmt.Println("get data : ", a)
	case b := <-strch:
		fmt.Println("get data : ", b)
	default:
		fmt.Println("no channel actvie")

	}

}



