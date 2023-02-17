package main

import (
	"fmt"
)

//交叉编译
//set CGO_ENABLED=0
//set GOOS=linux
//go env -w GOOS=linux
//go env -w GOOS=windows
//set GOARCH=amd64

func main() {
	fmt.Println(111)
}
