//捕获指定stack信息,一般在处理panic/recover中处理
//返回完整的堆栈信息和函数调用信息
package main

import (
	"bytes"
	"fmt"
	"log"
	"runtime"
)

func main() {
	test()
}

func test() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("recover error: %+v", err)

			// fmt.Fprintf(os.Stderr, "full stack info: \n%s", Stack())

			// stack := Stack()
			stack := CatchStack()
			log.Println("stack info: ", string(stack))
		}

	}()

	for i := 0; i < 10; i++ {
		if i == 5 {
			panic("exec panic")
		}

		log.Println("current index : ", i)
	}
}

//获取完整的堆栈信息
// Stack returns a formatted stack trace of the goroutine that calls it.
// It calls runtime.Stack with a large enough buffer to capture the entire trace.
func Stack() []byte {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			return buf[:n]
		}

		buf = make([]byte, 2*len(buf))
	}
}

//捕获指定stack信息,一般在处理panic/recover中处理
//返回完整的堆栈信息和函数调用信息
func CatchStack() []byte {
	buf := &bytes.Buffer{}

	//完整的堆栈信息
	stack := Stack()
	buf.WriteString("full stack:\n")
	buf.WriteString(string(stack))

	//完整的函数调用信息
	buf.WriteString("full fn call info:\n")

	// skip为0时，打印当前调用文件及行数。
	// 为1时，打印上级调用的文件及行数，依次类推
	for i := 1; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		fn := runtime.FuncForPC(pc).Name()
		buf.WriteString(fmt.Sprintf("error Stack file: %s:%d call func:%s\n", file, line, fn))
	}

	return buf.Bytes()
}