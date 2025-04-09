package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	// 使用通道替代队列，通道本身是线程安全的
	queue := make(chan int, 10)

	// 使用通道作为广播机制
	broadcastCh := make(chan struct{})

	// 使用原子变量跟踪队列状态
	var itemsProduced atomic.Int64

	// 创建等待组
	var wg sync.WaitGroup

	// 启动3个消费者
	for i := range 3 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			consumed := 0
			for consumed < 3 {
				select {
				case value, ok := <-queue:
					if !ok {
						return
					}
					consumed++
					fmt.Printf("消费者 %d: 消费数据 %d\n", id, value)
					// 模拟处理数据
					time.Sleep(time.Millisecond * 500)
				case <-broadcastCh:
					// 收到广播信号，可以执行特殊处理
					fmt.Printf("消费者 %d: 收到广播信号\n", id)
				}
			}
		}(i)
	}

	// 生产者
	go func() {
		for i := range 10 {
			queue <- i
			newCount := itemsProduced.Add(1)
			fmt.Printf("生产者: 生产数据 %d，已生产 %d 个元素\n", i, newCount)

			// 控制生产速度
			time.Sleep(time.Millisecond * 500)
		}
	}()

	// 等待一段时间后广播
	time.Sleep(time.Second * 5)
	fmt.Println("生产者: 广播通知所有消费者")
	close(broadcastCh) // 通过关闭通道实现广播

	// 确保所有数据被处理完
	time.Sleep(time.Second * 2)
	close(queue)

	// 等待所有消费者完成
	wg.Wait()
}
