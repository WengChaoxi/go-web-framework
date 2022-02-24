package framework

import "net/http"

// 框架核心数据结构
type Core struct {
	routers map[string]Handler
}

// 初始化框架核心结构
func NewCore() *Core {
	return &Core{
		routers: map[string]Handler{},
	}
}

func (c *Core) Get(url string, handler Handler) {
	c.routers[url] = handler
}

// 框架核心结构实现 Handler 接口
func (c *Core) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ctx := NewContext(req, rw)

	// test handler
	router := c.routers["test"]
	if router == nil {
		return
	}
	router(ctx)
}
