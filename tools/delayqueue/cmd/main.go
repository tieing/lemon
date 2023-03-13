package main

import (
	"context"
	"fmt"
	"github.com/tieing/lemon/tools/delayqueue"
	"time"
)

func consume(entry any) {
	fmt.Println("当前：", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("消费内容", entry.(string))
	fmt.Println("=======================")
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	q := delayqueue.New(ctx, delayqueue.WithFrequency(time.Millisecond*500))
	q.AddAfter(time.Millisecond*600, "600毫秒后")
	q.AddAfter(time.Second*15, "15秒后")
	q.AddAfter(time.Second*8, "8秒后")
	q.AddAfter(time.Second*43, "43秒后")
	q.AddAfter(time.Second*50, "50秒后")
	q.AddAfter(time.Second*28, "28秒后")

	q.Run()

	go func() {
		for {
			data, ok := q.Get()
			if !ok {
				println("关闭队列")
				return
			}
			consume(data)
		}
	}()
	select {
	case <-ctx.Done():
		println("exit")
	}
}
