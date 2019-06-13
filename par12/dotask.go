package main

/**
通过done chan信号执行goroutine并具有返回结果
*/
import (
	"log"
	"math/rand"
	"time"

	"github.com/daheige/thinkgo/common"
)

//产生[0-n]之间的[]int
func shuffle(n int) []int {
	rand.Seed(time.Now().UnixNano())
	b := rand.Perm(n)
	return b
}

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

	var x = 100 * 10000 //产生100w个数字打乱顺序
	res2 := common.DoTaskWithArgs(func(args ...interface{}) interface{} {
		n := args[0].(int)

		log.Println(shuffle(n))

		return nil
	}, x)

	//错误信息都可以抛出来
	log.Println("cost time: ", res2.CostTime) //cost time:  0.382877745
	log.Println(res2.Err, <-res2.Result)
}
