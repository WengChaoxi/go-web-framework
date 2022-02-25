package framework

import (
	"net/http"
	"strings"
)

// 框架核心数据结构
type Core struct {
	routers map[string]map[string]Handler
}

// 初始化框架核心结构
func NewCore() *Core {
	routers := map[string]map[string]Handler{}
	routers["GET"] = map[string]Handler{}
	routers["POST"] = map[string]Handler{}
	routers["PUT"] = map[string]Handler{}
	routers["DELETE"] = map[string]Handler{}
	return &Core{
		routers: routers,
	}
}

// Get
func (c *Core) Get(url string, handler Handler) {
	c.routers["GET"][strings.ToUpper(url)] = handler
}

// Post
func (c *Core) Post(url string, handler Handler) {
	c.routers["POST"][strings.ToUpper(url)] = handler
}

// Put
func (c *Core) Put(url string, handler Handler) {
	c.routers["PUT"][strings.ToUpper(url)] = handler
}

// Delete
func (c *Core) Delete(url string, handler Handler) {
	c.routers["DELETE"][strings.ToUpper(url)] = handler
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
		if handler, ok := methodHandlers[strings.ToUpper(uri)]; ok {
			return handler
		}
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
