package main

import (
	"fmt"
	"unsafe"
)

type test struct {
	len int
	str string
}

func main() {
	t1 := test{

		str: "123456",
	}
	t1.len = int(unsafe.Sizeof(t1.len) + unsafe.Sizeof(t1.str))
	fmt.Println(len(t1.str), unsafe.Sizeof(t1))
}
