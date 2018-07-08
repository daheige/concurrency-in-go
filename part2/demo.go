package main

import "sync"
import "fmt"

type Counter struct {
	mu    sync.Mutex
	value int
}

//利用mutex枷锁和解锁机制
// 通过使用内存访问同步原语，你可以隐藏从呼叫者锁定关键部分的实现细节，但不会给调用者带来复杂性。
// 这是一个线程安全类型的小例子：
func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func main() {
	count := &Counter{
		value: 1,
	}

	count.Increment()
	fmt.Println(count.value)
}
