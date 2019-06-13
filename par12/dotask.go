package main

/**
通过done chan信号执行goroutine并具有返回结果
*/
import (
	"log"

	"github.com/daheige/thinkgo/common"
)

func main() {
	//安全的执行goroutine
	res := common.DoTask(func() interface{} {
		for i := 0; i < 10000; i++ {
			log.Println("current index: ", i)
		}

		//这里故意抛出异常
		//panic(1) //当发生了panic后函数直接退出，不执行return
		return 1
	})

	log.Println("res: ", res)
	log.Println("exec time: ", res.CostTime)
	log.Println(res.Err == nil, <-res.Result == nil)
}
