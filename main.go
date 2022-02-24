package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"goweb/framework"
)

// handler func
func TestHandler(c *framework.Context) error {
	// id := c.QueryInt("id", 0)
	// if id == 1 {
	// 	c.Json(200, "hello world")
	// }
	// c.Json(200, "test")
	// return nil

	finish := make(chan struct{}, 1)
	panicChan := make(chan interface{}, 1)

	// 设置超时
	durationCtx, cancel := context.WithTimeout(c.BaseContext(), time.Duration(1*time.Second))
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

		// 业务逻辑
		id := c.QueryInt("id", 0)
		if id == 1 {
			c.Json(200, "hello world")
		}
		time.Sleep(2 * time.Second)
		c.Json(200, "test")

		finish <- struct{}{}
	}()

	select {
	case p := <-panicChan:
		c.WriterMux().Lock()
		defer c.WriterMux().Unlock()
		log.Println(p)
	case <-finish:
		fmt.Println("finish")
	case <-durationCtx.Done():
		c.WriterMux().Lock()
		defer c.WriterMux().Unlock()

		c.Json(500, "time out")
		c.SetHasTimeout()
	}
	return nil
}

func main() {
	core := framework.NewCore()
	core.Get("test", TestHandler)
	server := &http.Server{
		Handler: core,
		Addr:    ":80",
	}
	server.ListenAndServe()
}
