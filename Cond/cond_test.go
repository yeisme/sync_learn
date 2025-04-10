package cond

import (
	"sync"
	"testing"
)

func cond_run() {
	// 创建一个队列
	queue := make([]int, 0, 10)

	// 创建互斥锁和条件变量
	// 条件变量需要一个互斥锁来保护共享资源
	mutex := &sync.Mutex{}
	cond := sync.NewCond(mutex)

	// 创建一个等待组，用于等待所有消费者完成
	var wg sync.WaitGroup

	// 启动3个消费者
	for i := range 3 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for range 3 { // 每个消费者消费3个元素
				// 获取锁
				mutex.Lock()

				// 当队列为空时，等待条件变量
				for len(queue) == 0 {
					// fmt.Printf("消费者 %d: 队列为空，等待数据...\n", id)
					cond.Wait() // Wait会释放锁，并在被唤醒时重新获取锁
				}

				// 从队列取出数据
				_ = queue[0]
				queue = queue[1:]

				// fmt.Printf("消费者 %d: 消费数据 %d，队列剩余 %d 个元素\n", id, value, len(queue))

				// 释放锁
				mutex.Unlock()

				// 模拟处理数据
				// time.Sleep(time.Millisecond * 500)
			}
		}(i)
	}

	// 生产者
	go func() {
		for i := range 10 {
			// 获取锁
			mutex.Lock()

			// 向队列添加数据
			queue = append(queue, i)
			// fmt.Printf("生产者: 生产数据 %d，队列现有 %d 个元素\n", i, len(queue))

			// 唤醒一个等待的消费者
			cond.Signal()

			// 释放锁
			mutex.Unlock()

			// 控制生产速度
			// time.Sleep(time.Millisecond * 500)
		}
	}()

	// 等待一段时间后广播唤醒所有消费者
	// time.Sleep(time.Second * 5)
	mutex.Lock()
	// fmt.Println("生产者: 广播通知所有消费者")
	cond.Broadcast()
	mutex.Unlock()

	// 等待所有消费者完成
	wg.Wait()
}

func BenchmarkCond(b *testing.B) {
	for range b.N {
		cond_run()
	}
}
