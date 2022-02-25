package framework

type IGroup interface {
	Get(string, Handler)
	Post(string, Handler)
	Put(string, Handler)
	Delete(string, Handler)
	Group(string) IGroup // Group 嵌套
}

type Group struct {
	core   *Core  // 指向core，用于调用相关HTTP方法
	parent *Group // 指向上一个 Group，用于 Group 嵌套
	prefix string // 当前 Group 前缀
}

// 初始化 Group
func NewGroup(core *Core, prefix string) *Group {
	return &Group{
		core:   core,
		parent: nil,
		prefix: prefix,
	}
}

// Get
func (g *Group) Get(uri string, handler Handler) {
	uri = g.getAbsolutePrefix() + uri // 组合前缀和目标地址
	g.core.Get(uri, handler)          // 路由处理
}

// Post
func (g *Group) Post(uri string, handler Handler) {
	uri = g.getAbsolutePrefix() + uri
	g.core.Post(uri, handler)
}

// Put
func (g *Group) Put(uri string, handler Handler) {
	uri = g.getAbsolutePrefix() + uri
	g.core.Put(uri, handler)
}

// Delete
func (g *Group) Delete(uri string, handler Handler) {
	uri = g.getAbsolutePrefix() + uri
	g.core.Delete(uri, handler)
}

// 获取当前 Group 绝对路径
func (g *Group) getAbsolutePrefix() string {
	if g.parent == nil {
		return g.prefix
	}
	return g.parent.getAbsolutePrefix() + g.prefix
}

func (g *Group) Group(uri string) IGroup {
	group := NewGroup(g.core, uri)
	group.parent = g
	return group
}
