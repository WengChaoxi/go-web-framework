package framework

import (
	"log"
	"net/http"
	"strings"
)

// 框架核心数据结构
type Core struct {
	routers map[string]*Tree
}

// 初始化框架核心结构
func NewCore() *Core {
	routers := map[string]*Tree{}
	routers["GET"] = NewTree()
	routers["POST"] = NewTree()
	routers["PUT"] = NewTree()
	routers["DELETE"] = NewTree()
	return &Core{
		routers: routers,
	}
}

// Get
func (c *Core) Get(url string, handler Handler) {
	err := c.routers["GET"].AddRouter(url, handler)
	if err != nil {
		log.Fatal("add router error: ", err)
	} else {
		log.Println("add success")
	}
}

// Post
func (c *Core) Post(url string, handler Handler) {
	err := c.routers["POST"].AddRouter(url, handler)
	if err != nil {
		log.Fatal("add router error: ", err)
	}
}

// Put
func (c *Core) Put(url string, handler Handler) {
	err := c.routers["PUT"].AddRouter(url, handler)
	if err != nil {
		log.Fatal("add router error: ", err)
	}
}

// Delete
func (c *Core) Delete(url string, handler Handler) {
	err := c.routers["DELETE"].AddRouter(url, handler)
	if err != nil {
		log.Fatal("add router error: ", err)
	}
}

// Group
func (c *Core) Group(prefix string) IGroup {
	return NewGroup(c, prefix)
}

// 根据请求中的信息获取处理函数句柄
func (c *Core) FindRouterByRequest(request *http.Request) Handler {
	uri := request.URL.Path
	method := request.Method
	if methodHandlers, ok := c.routers[strings.ToUpper(method)]; ok {
		return methodHandlers.FindHandler(uri)
	}
	return nil
}

// 框架核心结构实现 Handler 接口
func (c *Core) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ctx := NewContext(req, rw)

	router := c.FindRouterByRequest(req)
	if router == nil {
		ctx.Json(404, "not found")
		return
	}
	if err := router(ctx); err != nil {
		ctx.Json(500, "inner error")
		return
	}
}
