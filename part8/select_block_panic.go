package main

import (
	"log"
	"time"
)

func main() {
	done := make(chan struct{}, 1)
	go func() {
		defer close(done)

		log.Println(111)
	}()

	<-done

	time.Sleep(1 * time.Second)
	select {} // go version >=14.x.x版本后，这里会发生携程阻塞，死锁，从而panic
}

/**
% go run app.go
2020/11/04 22:59:42 111
fatal error: all goroutines are asleep - deadlock!

go version >= 14.x.x版本后会发生panic，携程阻塞
goroutine 1 [select (no cases)]:
main.main()
        /Users/heige/web/go/demo/cmd/app.go:19 +0x85
*/
