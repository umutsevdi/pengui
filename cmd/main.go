package main

import (
	"fmt"
	"umutsevdi/pengui/sys"
)

func main() {
	fmt.Println(sys.GetHostInfo())
	for {
		fmt.Println(sys.Capture())
	}
}
