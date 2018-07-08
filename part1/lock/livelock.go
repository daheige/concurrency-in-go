package main

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	cadence := sync.NewCond(&sync.Mutex{})
	go func() {
		for range time.Tick(1 * time.Millisecond) {
			cadence.Broadcast() //广播
		}
	}()

	//takeStep模拟所有动作之间的恒定节奏
	takeStep := func() {
		cadence.L.Lock()
		cadence.Wait()
		cadence.L.Unlock()
	}

	tryDir := func(dirName string, dir *int32, out *bytes.Buffer) bool { //1
		fmt.Fprintf(out, " %v", dirName)
		atomic.AddInt32(dir, 1) //2
		takeStep()              //3
		if atomic.LoadInt32(dir) == 1 {
			fmt.Fprint(out, ". Success!")
			return true
		}
		takeStep()
		atomic.AddInt32(dir, -1) //4
		return false
	}

	var left, right int32
	tryLeft := func(out *bytes.Buffer) bool { return tryDir("left", &left, out) }
	tryRight := func(out *bytes.Buffer) bool { return tryDir("right", &right, out) }

	//移动步子
	walk := func(walking *sync.WaitGroup, name string) {
		var out bytes.Buffer
		defer func() { fmt.Println(out.String()) }()
		defer walking.Done()
		fmt.Fprintf(&out, "%v is trying to scoot:", name)
		for i := 0; i < 5; i++ { //1
			if tryLeft(&out) || tryRight(&out) { //2
				return
			}
		}
		fmt.Fprintf(&out, "\n%v tosses her hands up in exasperation!", name)
	}

	var peopleInHallway sync.WaitGroup //3 这个变量为程序提供了等待，直到两个人都能够相互通过或放弃
	peopleInHallway.Add(2)
	go walk(&peopleInHallway, "Alice")
	go walk(&peopleInHallway, "Barbara")
	peopleInHallway.Wait()

}

/**
tryDir 允许一个人尝试向某个方向移动并返回，无论他们是否成功。 每个方向都表示为试图朝这个方向移动的次数。
首先，我们通过将该方向递增1来朝着某个方向移动。 我们将在后面详细讨论atomic包。
现在，你只需要知道这个atomic包的操作是原子操作。
每个人必须以相同的速度或节奏移动。 takeStep模拟所有动作之间的恒定节奏。
在这里，这个人意识到他们不能在这个方向上放弃。 我们通过将该方向递减1来表示这一点
*/

/**运行结果
Barbara is trying to scoot: left right left right left right left right left right
Barbara tosses her hands up in exasperation!
Alice is trying to scoot: left right left right left right left right left right
Alice tosses her hands up in exasperation!
**/
/**
1. 我对尝试次数进行了人为限制，以便该程序结束。 在一个有活锁的程序中，可能没有这种限制，这就是为什么它是一个现实工作中的问题。
2. 首先，这个人会试图向左走，如果失败了，会尝试向右走。
3. 这个变量为程序提供了等待，直到两个人都能够相互通过或放弃
*/
/**
你可以看到Alice和Barbara在最终放弃之前持续交互。
这个例子演示了一个非常常见的活锁写入原因：
    两个或多个并发进程试图在没有协调的情况下防止死锁。 如果走廊里的人们一致认为只有一个人会移动
    那么就不会有活锁：一个人静止不动，另一个人移动到另一边，他们会继续走路。
*/
