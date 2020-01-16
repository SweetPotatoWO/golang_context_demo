package main

import (
	"context"
	"fmt"
	"time"
)

func main() {

	//一个函数  传递进去上下文 该函数的目的 就是无限的循环等待
	//并等待这context.Background 的关闭 也就是主协程
	gen := func(ctx context.Context) <-chan int {
		dst := make(chan int)
		n := 1
		go func() {
			for {
				select {
				case <-ctx.Done():
					fmt.Println("我老爸被关闭了,我也撤了")
					return
				case dst <- n:
					fmt.Println("加一")
					n++
				}
			}
		}()
		return dst
	}

	//withCancel 表示从 根父类中 继承和传递一个子的 context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() //调用撤销的方法

	//无限的获取 并无限的输出  知道n 为5的 时候 跳出循环
	//此时 主协程结束 并 cancel  context.Background 返回的父类context
	//通知 子类  ctx
	for n := range gen(ctx) {
		fmt.Println(n)
		if n == 100 {
			break
		}
	}
        //需要加个时间 延迟下 不然主程退出太过输出不了 Done()之后的字段
	time.Sleep(1 * time.Second)
}
