package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(1 * time.Second)

	//NewTicker+for实现每隔多久，指定一些操作
	i := 0
	for {
		<-ticker.C
		i++
		fmt.Println("i = ", i)

		if i == 5 {
			ticker.Stop() //停止定时器
			break
		}
	}
}
