//执行go build编译该demo
package main

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	msgCh := make(chan string, 6)
	done := make(chan struct{}, 1)

	//生产者和消费者运行模式： production:consumer = n:1
	//其他运行模式: 1:1 1:n n:n
	for i := 0; i < 10; i++ {
		go production(msgCh, done)
	}

	go consumer(msgCh)

	ch := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// recivie signal to exit main goroutine
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2, os.Interrupt, syscall.SIGHUP)

	// Block until we receive our signal.
	sig := <-ch
	log.Println("recive sig: ", sig.String())

	close(done) //when done closed,production will exit

	time.Sleep(3 * time.Second) //wait 3s

	log.Println("main goroutine will exit...")
}

func production(msgCh chan<- string, done <-chan struct{}) {
	i := 0
	for {
		select {
		case <-done:
			log.Println("production will exit...")
			return
		default:
			if i >= 1e6 {
				i = 0 //reset i
			}
		}

		rand.Seed(time.Now().UnixNano())
		if rand.Int31n(100) < 10 {
			continue
		}

		i++
		msgCh <- "hello,world: " + strconv.Itoa(i)
		log.Println("send data success!")
	}
}

func consumer(msgCh <-chan string) {
	for msg := range msgCh {
		log.Println("recive data: ", msg)
	}

}
