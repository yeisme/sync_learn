package main

import "sync"

func main() {

	once := sync.Once{}
	once.Do(func() {
		println("once.Do()")
	})

	(&sync.Once{}).Do(func() {
		println("once.Do()")
	})

	sync.OnceFunc(func() {
		println("once.Do()")
	})

}
