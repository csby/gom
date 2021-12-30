package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("info: http started at: ':8083'")
	defer fmt.Println("info: http stopped")

	err := http.ListenAndServe(":8083", &handler{})
	if err != nil {
		fmt.Println("error:", err)
	}
}
