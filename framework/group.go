package framework

type IGroup interface {
	Get(string, ...HandlerFunc)
	Post(string, ...HandlerFunc)
	Put(string, ...HandlerFunc)
	Delete(string, ...HandlerFunc)

	Group(string) IGroup // 用于 Group 嵌套

	Use(middlewares ...HandlerFunc)
}

type Group struct {
	core   *Core  // 指向core，用于调用相关HTTP方法
	parent *Group // 指向上一个 Group，用于 Group 嵌套
	prefix string // 当前 Group 前缀

	middlewares []HandlerFunc
}

// 初始化 Group
func NewGroup(core *Core, prefix string) *Group {
	return &Group{
		core:        core,
		parent:      nil,
		prefix:      prefix,
		middlewares: []HandlerFunc{},
	}
}

// Get
func (g *Group) Get(uri string, handlers ...HandlerFunc) {
	uri = g.getAbsolutePrefix() + uri // 组合前缀和目标地址
	allHandlers := append(g.getMiddlewares(), handlers...)
	g.core.Get(uri, allHandlers...)
}

// Post
func (g *Group) Post(uri string, handlers ...HandlerFunc) {
	uri = g.getAbsolutePrefix() + uri
	allHandlers := append(g.getMiddlewares(), handlers...)
	g.core.Get(uri, allHandlers...)
}

// Put
func (g *Group) Put(uri string, handlers ...HandlerFunc) {
	uri = g.getAbsolutePrefix() + uri
	allHandlers := append(g.getMiddlewares(), handlers...)
	g.core.Get(uri, allHandlers...)
}

// Delete
func (g *Group) Delete(uri string, handlers ...HandlerFunc) {
	uri = g.getAbsolutePrefix() + uri
	allHandlers := append(g.getMiddlewares(), handlers...)
	g.core.Get(uri, allHandlers...)
}

// 获取当前 Group 绝对路径
func (g *Group) getAbsolutePrefix() string {
	if g.parent == nil {
		return g.prefix
	}
	return g.parent.getAbsolutePrefix() + g.prefix
}

// 获取当前 Group 的 middleware
func (g *Group) getMiddlewares() []HandlerFunc {
	if g.parent == nil {
		return g.middlewares
	}
	return append(g.parent.getMiddlewares(), g.middlewares...)
}

func (g *Group) Group(uri string) IGroup {
	group := NewGroup(g.core, uri)
	group.parent = g
	return group
}

// 注册中间件
func (g *Group) Use(middlewares ...HandlerFunc) {
	g.middlewares = append(g.middlewares, middlewares...)
}
