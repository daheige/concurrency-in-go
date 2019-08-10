package main

import (
	"log"
	"net/http"

	"github.com/didip/tollbooth"
)

func HelloHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func main() {
	// Create a request limiter per handler.
	//每秒100个
	limiter := tollbooth.NewLimiter(100, nil)

	http.Handle("/", tollbooth.LimitFuncHandler(limiter, HelloHandler))

	address := ":4000"
	log.Println("server has run: ", address)
	http.ListenAndServe(address, nil)
}
