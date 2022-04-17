package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/WengChaoxi/go-web-framework/framework"
	"github.com/WengChaoxi/go-web-framework/framework/middleware"
)

// test
func TestHandler(c *framework.Context) error {
	id := c.QueryInt("id", 0)
	if id == 1 {
		c.Json(200, "hello world")
	} else if id == 2 {
		c.Text("test")
	} else if id == 3 {
		c.HTML("<h1>你好世界</h1>", nil)
	}
	return nil
}

// default
func DefaultHandler(c *framework.Context) error {
	c.Json(200, c.Request().URL.Path)
	return nil
}

func main() {
	core := framework.NewCore()

	// 超时中间件
	core.Use(middleware.Timeout(1 * time.Second))

	// 注册 /test 路由
	core.Get("/test", TestHandler)

	// 使用 group 分组，路由前缀
	userRouter := core.Group("/user")
	{
		userRouter.Use(middleware.Cost()) // 为当前分组使用计算耗时中间件
		userRouter.Get("/login", DefaultHandler)
		userRouter.Get("/logout", DefaultHandler)
	}

	server := &http.Server{
		Handler: core,
		Addr:    ":80",
	}

	// 启动 http 服务
	go func() {
		server.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)

	// ctrl+c  : SIGINT
	// ctrl+\  : SIGQUIT
	// kill    : SIGTERM
	// kill -9 : SIGKILL // 不能被捕获
	signal.Notify(quit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	<-quit

	fmt.Println("shutdown...")
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	<-ticker.C

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(timeoutCtx); err != nil {
		log.Fatal("server shutdown: ", err)
	}
}
