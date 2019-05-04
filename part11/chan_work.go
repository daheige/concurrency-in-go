package main

import "fmt"

func Task(taskch, resch chan int, exitch chan bool) {
	defer func() {   //异常处理
		err := recover()
		if err != nil {
			fmt.Println("do task error：", err)
			return
		}
	}()

	for t := range taskch { //  处理任务
		fmt.Println("do task :", t)
		resch <- t //
	}
	exitch <- true //处理完发送退出信号
}

func main() {
	taskch := make(chan int, 20) //任务管道
	resch := make(chan int, 20)  //结果管道
	exitch := make(chan bool, 5) //退出管道
	go func() {
		for i := 0; i < 10; i++ {
			taskch <- i
		}
		close(taskch)
	}()


	for i := 0; i < 5; i++ {  //启动5个goroutine做任务
		go Task(taskch, resch, exitch)
	}

	go func() { //等5个goroutine结束
		for i := 0; i < 5; i++ {
			<-exitch
		}

		close(resch)  //任务处理完成关闭结果管道，不然range报错
		close(exitch)  //关闭退出管道
	}()

	//读取结果
	for res := range resch{  //打印结果
		fmt.Println("task res：",res)
	}
}

