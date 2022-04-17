package framework

import (
	"net/http"
	"strings"
)

// 框架核心数据结构
type Core struct {
	routers     map[string]map[string][]HandlerFunc
	middlewares []HandlerFunc
}

// 初始化框架核心结构
func NewCore() *Core {
	routers := map[string]map[string][]HandlerFunc{}
	routers["GET"] = map[string][]HandlerFunc{}
	routers["POST"] = map[string][]HandlerFunc{}
	routers["PUT"] = map[string][]HandlerFunc{}
	routers["DELETE"] = map[string][]HandlerFunc{}
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
	c.routers["GET"][url] = allHandlers
}

// Post
func (c *Core) Post(url string, handlers ...HandlerFunc) {
	allHandlers := append(c.middlewares, handlers...)
	c.routers["POST"][url] = allHandlers
}

// Put
func (c *Core) Put(url string, handlers ...HandlerFunc) {
	allHandlers := append(c.middlewares, handlers...)
	c.routers["PUT"][url] = allHandlers
}

// Delete
func (c *Core) Delete(url string, handlers ...HandlerFunc) {
	allHandlers := append(c.middlewares, handlers...)
	c.routers["DELETE"][url] = allHandlers
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
		return methodHandlers[uri]
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
	ctx.Next()
}
