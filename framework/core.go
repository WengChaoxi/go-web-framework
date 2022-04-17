package framework

import (
	"log"
	"net/http"
	"strings"
)

// 框架核心数据结构
type Core struct {
	routers     map[string]*Tree
	middlewares []HandlerFunc
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

// 注册中间件
func (c *Core) Use(middlewares ...HandlerFunc) {
	c.middlewares = append(c.middlewares, middlewares...)
}

// Get
func (c *Core) Get(url string, handlers ...HandlerFunc) {
	allHandlers := append(c.middlewares, handlers...)
	err := c.routers["GET"].AddRouter(url, allHandlers)
	if err != nil {
		log.Fatal("add router error: ", err)
	}
}

// Post
func (c *Core) Post(url string, handlers ...HandlerFunc) {
	allHandlers := append(c.middlewares, handlers...)
	err := c.routers["POST"].AddRouter(url, allHandlers)
	if err != nil {
		log.Fatal("add router error: ", err)
	}
}

// Put
func (c *Core) Put(url string, handlers ...HandlerFunc) {
	allHandlers := append(c.middlewares, handlers...)
	err := c.routers["PUT"].AddRouter(url, allHandlers)
	if err != nil {
		log.Fatal("add router error: ", err)
	}
}

// Delete
func (c *Core) Delete(url string, handlers ...HandlerFunc) {
	allHandlers := append(c.middlewares, handlers...)
	err := c.routers["DELETE"].AddRouter(url, allHandlers)
	if err != nil {
		log.Fatal("add router error: ", err)
	}
}

// Group
func (c *Core) Group(prefix string) IGroup {
	return NewGroup(c, prefix)
}

// 根据请求中的信息获取路由对应处理函数句柄
func (c *Core) FindHandlersByRequest(request *http.Request) []HandlerFunc {
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

	handlers := c.FindHandlersByRequest(req)
	if handlers == nil {
		ctx.Json(404, "not found")
		return
	}
	ctx.SetHandlers(handlers)

	// if err := ctx.Next(); err != nil {
	// 	ctx.Json(500, "inner error")
	// 	return
	// }
	ctx.Next()
}
