package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	c := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		close(c) // 1
	}()

	fmt.Println("Blocking on read...")

	//当关闭了通道后,可以读取通道的默认值
	select {
	case v := <-c: // 2
		fmt.Println(v)
		fmt.Printf("Unblocked %v later.\n", time.Since(start))
	}
}

/**运行结果
Blocking on read...
<nil>
Unblocked 5.000333036s later.
*/
