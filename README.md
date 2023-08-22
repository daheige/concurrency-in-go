# concurrency-in-go

《Go 语言并发之道》笔记

如下部分，用实际的例子分析了 go 并发编程为什么难写，以及对 go 并发编程进行了全面剖析。

- [part1](./part1/): 并发编程为什么难写
- part2: 代码建模:序列化交互处理
- part3: go 的并发构建模块方案和 goroutine 通信方式
- part4: go 的并发编程范式
- part5: 可伸缩的并发设计
- part6: chan+goroutine 使用和 go 运行时任务调度
- part7: time.Ticker,Time.Tick 定时器用法
- part8: select 用法
- part9: 关于 goroutine 出现惊慌或致命错误的处理
- part10: 关于 chan 在通信过程中，生产者和消费者运行模式研究
- part11: goroutine 调度机制 MPG
- part12: 采用 context+select 实现超时调用
- part13: golang 多核 cpu 计算
- part14: golang interface 实战
- part15: 控制 goroutine 并发数
- part16: 限流控制
- part17: go 程序平滑退出机制
- part18: 关于堆栈和逃逸分析
- part19: 关于 sync.WaitGroup 协程计数器和 goroutine 通信机制
- part20: safe go 安全的 goroutine 执行
- part21: go work pool 探讨
- part22: go timeout 探讨和 panic 捕获处理
- part23: panic 在 http server 如何捕获

# 关于 http server main 平滑退出

对于 go v1.14.x 版本之后，goroutine 可抢夺，所以 select{}方式不可取

参考：
https://github.com/daheige/concurrency-in-go/blob/master/part8/select_block_exit.go

# 参考文档

- https://www.kancloud.cn/mutouzhang/go/596838
- https://github.com/daheige/Go-Questions/blob/master/channel/12%20-%20channel%20%E6%9C%89%E5%93%AA%E4%BA%9B%E5%BA%94%E7%94%A8.md

# License

    MIT
