package main

import (
	"fmt"
)

func fun(a int) (err string) {
	if a > 0 {
		err = "错误"
	}
	return
}

func main() {
	err := fun(5)
	fmt.Println(err)
	if err == "" {
		fmt.Println("return了正确的返回值")
	} else {
		fmt.Println("return了错误的返回值")
	}
}
