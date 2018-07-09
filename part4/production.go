package main

import "fmt"

//并发安全:将通道的读写约束在指定范围内
func main() {
	chanOwner := func() <-chan int {
		results := make(chan int, 5) //1
		go func() {
			defer close(results)
			for i := 0; i <= 5; i++ {
				fmt.Println("send data", i)
				results <- i
			}
		}()
		return results
	}

	consumer := func(results <-chan int) { //3
		for result := range results {
			fmt.Printf("Received: %d\n", result)
		}
		fmt.Println("Done receiving!")
	}

	results := chanOwner() //2
	consumer(results)
}
