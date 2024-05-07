package main

import (
	"fmt"

	"github.com/subinoybiswas/goenv"
)

func main() {
	val, err := goenv.GetEnv("HI")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(val)
}
