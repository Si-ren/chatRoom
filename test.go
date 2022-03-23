package main

import (
	"fmt"
	"unsafe"
)

func main() {
	var b int32 = 5
	fmt.Println(unsafe.Sizeof(b))
}
