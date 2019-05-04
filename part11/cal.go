package main

import (
	"fmt"
	"time"
)

func cal(a int , b int )  {
	c := a+b
	fmt.Printf("%d + %d = %d\n",a,b,c)
}

func main() {
	for i := 0; i < 10; i++ {
		go cal(i, i+1) //启动10个goroutine 来计算
	}

	Arry := make([]int,4) //4个元素长度
	for i :=0 ; i<10 ;i++{
		go addele(Arry,i)
	}

	time.Sleep(time.Second * 2) // sleep作用是为了等待所有任务完成
}

//当启动多个goroutine时，如果其中一个goroutine异常了，并且我们并没有对进行异常处理，那么整个程序都会终止
//所以我们在编写程序时候最好每个goroutine所运行的函数都做异常处理，异常处理采用recover
func addele(a []int ,i int)  {
	//异常捕获
	defer func() {    //匿名函数捕获错误
		err := recover()
		if err != nil {
			fmt.Println("add ele fail")
		}
	}()

	a[i]=i
	fmt.Println(a)
}