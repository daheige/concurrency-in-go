package main

import (
	"fmt"
)

//假设我们有一个数据竞争：两个并发进程试图访问同一个内存区域，并且它们访问内存的方式不是原子的
func main() {
	var data int
	go func() { data++ }()
	if data == 0 {
		fmt.Println("the value is 0.")
	} else { //保证数据有输出
		fmt.Printf("the value is %v.\n", data)
	}
}

/**
程序中有一些操作需要独占访问共享资源。在这个例子中，我们找到三处：

goroutine正在增加数据变量。
if语句，它检查数据的值是否为0。
fmt.Printf语句，用于检索输出数据的值。
有很多方法可以保护这些访问，Go有很好的方式来处理这个问题，解决这个问题的方法之一是让这些操作同步访问内存
*/
