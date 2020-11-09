# concurrency-in-go

    如下部分，用实际的例子分析了go并发编程为什么难写，以及对go chan,goroutine进行了全面剖析。
    part1: 并发编程为什么难写
    part2: 代码建模:序列化交互处理
    part3: go的并发构建模块方案和goroutine通信方式
    part4: go的并发编程范式
    part5: 可伸缩的并发设计
    part6: chan+goroutine使用和go运行时任务调度
    part7: time.Ticker,Time.Tick定时器用法
    part8: select用法
    part9: 关于goroutine出现惊慌或致命错误的处理
    part10: 关于chan在通信过程中，生产者和消费者运行模式研究
    part11: goroutine调度机制MPG
    part12: 采用context+select实现超时调用
    part13: golang多核cpu计算
    part14: golang interface实战
    part15: 控制goroutine并发数
    part16: 限流控制
    part17: go程序平滑退出机制
    part18: 关于堆栈和逃逸分析
    part19: 关于sync.WaitGroup协程计数器和goroutine通信机制
    part20: safe go安全的goroutine执行
    part21: go work pool探讨
    part22: goroutine panic捕获处理以及超时处理探讨

# 关于main退出
    对于go v1.14.x版本之后，goroutine可抢夺，所以select{}方式不可取
    参考：https://github.com/daheige/concurrency-in-go/blob/master/part8/select_block_exit.go
    
# 参考文档

    https://www.kancloud.cn/mutouzhang/go/596838
    https://github.com/daheige/Go-Questions/blob/master/channel/12%20-%20channel%20%E6%9C%89%E5%93%AA%E4%BA%9B%E5%BA%94%E7%94%A8.md

# License

    MIT
