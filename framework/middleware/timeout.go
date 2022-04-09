package middleware

import (
	"context"
	"fmt"
	"log"
	"time"

	"goweb/framework"
)

func Timeout(duration time.Duration) framework.HandlerFunc {
	return func(c *framework.Context) error {
		finish := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)

		// 设置超时
		durationCtx, cancel := context.WithTimeout(c.BaseContext(), duration)
		defer cancel()

		go func() {
			// 根据 golang 的设计，每个 Goroutine 都是独立存在的
			// 父 Goroutine 一旦使用 go 关键字开启一个子 Goroutine
			// 父子是平等的，将不能相互干扰，都必须自己处理自己的异常
			// 任意一个 Goroutine 的 panic 都会导致整个进程崩溃
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()

			// 处理业务逻辑
			c.Next()

			finish <- struct{}{}
		}()

		select {
		case p := <-panicChan:
			c.WriterMux().Lock()
			defer c.WriterMux().Unlock()
			log.Println(p)
			c.GetResponse().WriteHeader(500)
		case <-finish:
			fmt.Println("finish")
		case <-durationCtx.Done():
			c.WriterMux().Lock()
			defer c.WriterMux().Unlock()

			c.Json(504, "time out")
			c.SetHasTimeout()
		}
		return nil
	}
}
