//通过枷锁的方式保证数据的同步
//Go的惯用法（我不建议你像这样解决数据竞争问题），但它很简单地演示了内存访问同步
package main

import (
	"fmt"
	"sync"
)

var mutex sync.Mutex //1

func main() {
	var data int
	go func() {
		//对data的叠加,枷锁独占访问数据变量的内存
		mutex.Lock()
		data++
		mutex.Unlock() //操作完毕后释放锁
	}()

	//读取数据的时候,也枷锁和释放锁
	mutex.Lock()
	if data == 0 {
		fmt.Println("the value is 0.")
	} else { //保证数据有输出
		fmt.Printf("the value is %v.\n", data)
	}
	mutex.Unlock()
	fmt.Println("read data compleated")
}

/**
虽然我们已经解决了数据竞争，但我们并没有真正解决竞争条件！这个程序的操作顺序仍然不确定。
我们刚刚只是缩小了非确定性的范围。
在这个例子中，仍然不确定goroutine是否会先执行，或者我们的if和else块是否都会执行
它不会自动解决数据竞争或逻辑正确性问题,此外，它还可能导致维护和性能问题

性能分析;
以这种方式同步对内存的访问会导致性能下降。每次我们执行其中一项操作时，程序会暂停一段时间。 这带来了两个问题：
    加锁的程序部分是否重复进入和退出？
    加锁的程序对内存占用到底有多大？
    要说清这两个问题简直是门艺术,同步对内存的访问也与其他并发建模存在关联
**/
