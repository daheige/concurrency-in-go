package main

import (
	"log"
)

type User struct {
	Id   int
	Name string
}

// getUser 返回一个指针类型
func getUser() *User {
	log.Println("hello")
	return &User{Id: 1, Name: "daheige"}
}

// getUser2 模拟参数泄露
func getUser2(u *User) *User {
	log.Println("hello2")
	return u
}

// getUser3 避免参数泄露，让返回值分配在堆上
//.\app.go:25:15: moved to heap: u
func getUser3(u User) *User {
	log.Println("hello3")
	return &u
}

func main() {
	//_ = getUser()

	u := getUser()
	log.Println(u)

	str := new(string) //分配内存到了堆上
	*str = "daheige"
	//fmt.Println(str)

	//泄露参数
	getUser2(&User{Id: 1, Name: "daheige2"})

	u3 := getUser3(User{Id: 1, Name: "daheige3"})
	log.Println(u3) //.\app.go:44:13: u3 escapes to heap

}

/**
变量&User发生了逃逸到堆上,分配到了堆上
>go build -gcflags "-m -l" -o app
# part1
.\app.go:12:13: getUser ... argument does not escape
.\app.go:12:14: "hello" escapes to heap
.\app.go:13:28: &User literal escapes to heap

通过汇编来看
>go tool compile -S app.go
	0x0073 00115 (app.go:13)        CALL    runtime.newobject(SB)
将目光集中到 CALL 指令，发现其执行了 runtime.newobject 方法，也就是确实是分配到了堆上

GetUserInfo() 返回的是指针对象，引用被返回到了方法之外了。
因此编译器会把该对象分配到堆上，而不是栈上。否则方法结束之后，局部变量就被回收了，岂不是翻车。所以最终分配到堆上是理所当然的

go build -gcflags "-m -l" -o app.exe
# part1
.\app.go:12:13: getUser ... argument does not escape
.\app.go:12:14: "hello" escapes to heap
.\app.go:13:28: &User literal escapes to heap
.\app.go:20:13: main ... argument does not escape
.\app.go:20:13: u escapes to heap

// 同样发生了逃逸
// str 变量逃到了堆上，也就是该对象在堆上分配
.\app.go:25:12: new(string) escapes to heap
.\app.go:28:13: main ... argument does not escape
.\app.go:28:13: str escapes to heap

案例二只加了一行代码 fmt.Println(str)，问题肯定出在它身上。其原型：

func Println(a ...interface{}) (n int, err error)
通过对其分析，可得知当形参为 interface 类型时，在编译阶段编译器无法确定其具体的类型。
因此会产生逃逸，最终分配到堆上
如果你有兴趣追源码的话，可以看下内部的 reflect.TypeOf(arg).Kind() 语句，其会造成堆逃逸，
而表象就是 interface 类型会导致该对象分配到堆上
当我们把fmt.Println这一行去掉，再次运行
$ go build -gcflags "-m -l" -o app.exe
.\app.go:24:12: main new(string) does not escape
str就不会发生逃逸

$ go build -gcflags "-m -l" -o app.exe
.\app.go:19:15: leaking param: u to result ~r1 level=0
发现这一行出现了 leaking param
它说明了变量 u 是一个泄露参数。结合代码可得知其传给 GetUserInfo 方法后，没有做任何引用之类的涉及变量的动作，
直接就把这个变量返回出去了。因此这个变量实际上并没有逃逸，它的作用域还在 main() 之中，所以分配在栈上
这个变量没发生逃逸，但是发生了参数泄露，变量还在当前上下文中，也就是main中，分配在栈上，当main退出，就会自动销毁

$ go build -gcflags "-m -l" -o app.exe
//.\app.go:44:13: u3 escapes to heap
当把参数u User传递给getUser3 ，参数没有发生逃逸，返回值&u发生变量逃逸到堆上

小结:
	静态分配到栈上，性能一定比动态分配到堆上好
	底层分配到堆，还是栈。实际上对你来说是透明的，不需要过度关心
	每个 Go 版本的逃逸分析都会有所不同（会改变，会优化）
	直接通过 go build -gcflags '-m -l' 就可以看到逃逸分析的过程和结果
	到处都用指针传递并不一定是最好的，要用对
*/
