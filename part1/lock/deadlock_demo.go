package main

import (
	"fmt"
	"sync"
	"time"
)

type value struct {
	mu    sync.Mutex
	value int
}

func main() {
	var wg sync.WaitGroup
	//尝试以枷锁和释放锁的方式访问变量v1,v2
	printSum := func(v1, v2 *value) {
		defer wg.Done()
		v1.mu.Lock()         //1 访问带锁的部分
		defer v1.mu.Unlock() //2 试图调用defer关键字释放锁

		time.Sleep(2 * time.Second) //3 添加休眠时间 以造成死锁
		v2.mu.Lock()
		defer v2.mu.Unlock()

		fmt.Printf("sum=%v\n", v1.value+v2.value)
	}

	var a, b value
	wg.Add(2)

	//让a,b相互等待
	go printSum(&a, &b)
	go printSum(&b, &a)
	wg.Wait()
}

//运行报错
// fatal error: all goroutines are asleep - deadlock!
// 实质上，我们创建了两个不能一起运转的齿轮：我们的第一个打印总和调用a锁定，然后尝试锁定b，但与此同时，我们打印总和的第二个调用锁定了b并尝试锁定a。
// 两个goroutine都无限地等待着彼此

/*出现死锁的分析(科夫曼条件分析法)
  1. printSum函数确实需要a和b的独占权，所以它满足了这个条件。
  2. 因为printSum保持a或b并等待另一个，所以它满足这个条件。
  3. 我们没有任何办法让我们的goroutine被抢占。
  4. 我们第一次调用printSum正在等待我们的第二次调用，反之亦然。
*/
