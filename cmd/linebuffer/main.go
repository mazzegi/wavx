package main

import (
	"fmt"

	"github.com/mazzegi/wavx"
)

func main() {
	b := wavx.NewLinearBuffer(5)
	for i := 0; i < 20; i++ {
		b.Push(float64(i))
		fmt.Println(b.FirstValue())
	}
}
