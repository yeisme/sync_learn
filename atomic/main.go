package main

import (
	"fmt"
	"sync/atomic"
)

func main() {
	a_p := new(atomic.Int32)
	a := atomic.Int32{}

	a_p.Store(1)

	if a_p.Load() != 1 {
		panic("expected 1")
	}

	fmt.Printf("atomic.Int32: %v\n", a_p)
	fmt.Printf("atomic.Int32: %v\n", a.Load())
}
