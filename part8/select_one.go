package main

import (
	"log"
	"time"
)

func main() {
	name := query([]Conn{
		Conn{
			name: "xiaoming",
		},
		Conn{
			name: "xiaoming",
		},
	})

	log.Println("name = ", name)
}

// Conn 模拟connection
type Conn struct {
	name string
}

// DoQuery 模拟查询
func (c Conn) DoQuery(n int) string {
	if n%2 == 0 {
		return c.name
	}

	time.Sleep(time.Duration(2*n+1) * time.Millisecond)

	return c.name
}

func query(conns []Conn) interface{} {
	ch := make(chan interface{}, 1)

	for k, conn := range conns {
		go func(c Conn) {
			select {
			case ch <- c.DoQuery(k):
				log.Println("current conn key: ", k)
				log.Println("query end")
			default:
			}

		}(conn)
	}

	return <-ch
}

/**
2020/05/10 17:35:23 current conn key:  1
2020/05/10 17:35:23 query end
2020/05/10 17:35:23 name =  xiaoming
*/
