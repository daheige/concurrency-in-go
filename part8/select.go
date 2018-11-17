package main

import "fmt"

func main() {
	x := makeRndNum()
	fmt.Println(x)
}

func makeRndNum() int {
	x := make(chan int)
	//独立goroutine生成0,1
	//当select中的case超过1后，就会形成阻塞模式，会按照select内部算法调用case执行
	//如果指定了default,当select没有获得通道操作权限，就会执行default操作
	go func() {
		select { //随机选择一个chan写入
		case x <- 1:
			fmt.Println("x = 1")
		case x <- 0:
			fmt.Println("x = 0")
		}
	}()

	return <-x
}
