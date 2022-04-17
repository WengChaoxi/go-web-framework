// 字典树，用于实现动态路由匹配
package framework

import (
	"errors"
	"strings"
)

type Tree struct {
	root *node
}

type node struct {
	isLast   bool          // 代表当前节点是否可以成为最终的路由规则
	segment  string        // uri 中的字符串，路由的片段
	handlers []HandlerFunc // 中间件、路由处理函数句柄
	childs   []*node       // 所属子节点
}

func newNode() *node {
	return &node{
		isLast:  false,
		segment: "",
		childs:  []*node{},
	}
}

func NewTree() *Tree {
	return &Tree{
		root: newNode(),
	}
}

// 判断一个 segment 中是否以 : 开头
func isWildSegment(segment string) bool {
	return strings.HasPrefix(segment, ":") // 以 : 开头则是通用 segment
}

// 过滤下一层满足 segment 规则的子节点
func (n *node) filterChildNodes(segment string) []*node {
	if len(n.childs) == 0 {
		return nil
	}

	if isWildSegment(segment) {
		return n.childs
	}

	nodes := make([]*node, 0, len(n.childs))
	// 过滤所有的下一层子节点
	for _, node := range n.childs {
		if isWildSegment(node.segment) { // 如果下一层子节点有通用路由
			nodes = append(nodes, node)
		} else if node.segment == segment { // 如果下一层子节点没有通用路由，但是文本完全匹配
			nodes = append(nodes, node)
		}
	}
	return nodes
}

// 查找
func (n *node) matchNode(uri string) *node {
	segments := strings.SplitN(uri, "/", 2)
	segment := segments[0]
	if !isWildSegment(segment) {
		segment = strings.ToUpper(segment)
	}
	nodes := n.filterChildNodes(segment)
	// 当前子节点没有符合条件的，则uri不存在
	if len(nodes) == 0 {
		return nil
	}
	// 只有1个 segment
	if len(segments) == 1 {
		for _, node := range nodes {
			if node.isLast {
				return node
			}
		}
		return nil
	}
	// 有2个 segment，递归查找
	for _, node := range nodes {
		match := node.matchNode(segments[1])
		if match != nil {
			return match
		}
	}

	return nil
}

// 添加路由节点
func (t *Tree) AddRouter(uri string, handlers []HandlerFunc) error {
	n := t.root
	if n.matchNode(uri) != nil {
		return errors.New("route exist: " + uri)
	}
	segments := strings.Split(uri, "/")
	for index, segment := range segments {
		if !isWildSegment(segment) {
			segment = strings.ToUpper(segment)
		}
		// 是否到最后节点
		isLast := index == len(segments)-1

		var tmp *node
		childNodes := n.filterChildNodes(segment)
		if len(childNodes) > 0 {
			for _, node := range childNodes {
				if node.segment == segment {
					tmp = node
					break
				}
			}
		}

		// 不存在，则创建
		if tmp == nil {
			node := newNode()
			node.segment = segment
			if isLast {
				node.isLast = true
				node.handlers = handlers
			}
			n.childs = append(n.childs, node)
			tmp = node
		}
		n = tmp
	}
	return nil
}

// 匹配 uri
func (t *Tree) FindHandler(uri string) []HandlerFunc {
	matchNode := t.root.matchNode(uri)
	if matchNode != nil {
		return matchNode.handlers
	}
	return nil
}
