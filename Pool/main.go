package main

import (
	"bytes"
	"fmt"
	"sync"
	"time"
)

// 自定义对象类型，模拟大型对象
type BigObject struct {
	buffer bytes.Buffer
	data   [1024]byte // 1KB 的数据
}

// 重置对象以便复用
func (o *BigObject) Reset() {
	o.buffer.Reset()
}

func main() {
	// 创建对象池
	pool := sync.Pool{
		// 函数签名需要遵循 sync.Pool 的接口定义，返回 interface{}/any 类型
		New: func() any {
			fmt.Println("创建新的 BigObject")
			return &BigObject{}
		},
	}

	// 模拟并发任务
	var wg sync.WaitGroup

	// 执行 100 个并发任务
	for i := range 100 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// 从池中获取对象
			// func (p *sync.Pool) Get() any
			obj := pool.Get().(*BigObject)
			// 确保函数结束时将对象放回池中
			defer pool.Put(obj)

			// 重置对象以便安全复用
			obj.Reset()

			// 使用对象进行一些操作
			obj.buffer.WriteString(fmt.Sprintf("任务 %d 正在处理数据", id))

			// 模拟一些处理时间
			time.Sleep(time.Millisecond)
		}(i)
	}

	wg.Wait()
	fmt.Println("所有任务已完成")

	// 展示池的复用特性
	for i := range 5 {
		obj := pool.Get().(*BigObject)
		fmt.Printf("# %d 获取到对象: %p\n", i, obj)
		// 不返回最后几个对象以观察新对象创建
		if i < 3 {
			pool.Put(obj)
		}
		/*
			# 0 获取到对象: 0xc000191b08
			# 1 获取到对象: 0xc000191b08
			# 2 获取到对象: 0xc000191b08
			# 3 获取到对象: 0xc000191b08
			# 4 获取到对象: 0xc00011cd88
		*/
	}
}
