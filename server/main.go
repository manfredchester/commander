package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", handle)
	http.ListenAndServe("127.0.0.1:1122", nil)
}

func handle(rsp http.ResponseWriter, req *http.Request) {
	fmt.Println("handler started!")
	defer fmt.Println("handler ended!")

	ctx := req.Context()

	select {
	case <-ctx.Done():
		fmt.Println("stop!")
		fmt.Println(ctx.Err())
		http.Error(rsp, ctx.Err().Error(), http.StatusInternalServerError)
	case <-time.After(5 * time.Second):
		fmt.Println("wait 5 seconds")
		fmt.Fprintln(rsp, "hello!")
	}
}
