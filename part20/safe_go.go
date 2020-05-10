package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	err := New(func() error {
		log.Println("111")
		arr := []string{
			"a",
			"b",
			"c",
		}

		// 故意模拟slice取值越界
		log.Println(arr[3])

		return nil
	}).Do()

	log.Println("err: ", err)
}

// logger 记录日志
type logger interface {
	Println(args ...interface{})
}

// New returns 创建一个safeGo实例
func New(fn func() error, logEntry ...logger) *safeGo {
	s := &safeGo{
		fn: fn,
	}

	if len(logEntry) > 0 {
		s.logEntry = logEntry[0]
	} else {
		s.logEntry = log.New(os.Stderr, "", log.LstdFlags)
	}

	return s
}

// safeGo safe go
type safeGo struct {
	fn       func() error
	logEntry logger
}

// Do 安全的执行fn
func (s safeGo) Do() (err error) {
	defer func() {
		if e := recover(); e != nil {
			s.logEntry.Println("current fn exec panic: ", e)
			err = fmt.Errorf("exec panic: %v", e)
		}
	}()

	err = s.fn()
	return
}
