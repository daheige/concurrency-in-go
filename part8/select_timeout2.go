package main

import (
	"errors"
	"log"
	"time"
)

func main() {
	ch := make(chan error, 1)

	// 在独立的协程中运行
	go func() {
		time.Sleep(2 * time.Second)
		log.Println("nock timeout")
		ch <- errors.New("do something timeout")
	}()

	select {
	case resp := <-ch:
		log.Println("res ", resp)
	case <-time.After(1 * time.Second):
		log.Println("wait timeout")
	}

}
