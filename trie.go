package ree

import (
	rlog "Ree/ree/logs"
	"errors"
	"strings"
)

type node struct {
	isWild   bool             // 是否是*
	children map[string]*node // 子节点
	part     string           // /分割的路径
	path     string           // 完整的路径
}

func parsePath(path string) []string {
	return strings.Split(path, "/")
}

func (r *Route) addRoute(method string, path string, handler HandlerFunc) error {
	parts := parsePath(path)[1:]
	if _, ok := r.trie[method]; !ok {
		r.trie[method] = &node{children: make(map[string]*node)}
	}
	n := r.trie[method]
	for i, part := range parts {
		_, ok := n.children[part]
		if ok && n.isWild {
			return errors.New("路由冲突")
		}
		if !ok {
			n.children[part] = &node{
				isWild:   part[0] == ':' || part[0] == '*',
				children: make(map[string]*node),
				part:     part,
				path:     strings.Join(parts[:i], "/"),
			}
		}
		n = n.children[part]
	}
	n.path = path
	n.children = nil
	r.handlers[method+"-"+path] = handler
	return nil
}

func (n *node) print(method string) {
	if len(n.children) == 0 {
		rlog.Info("%s--%s", method, n.path)
		return
	}

	for _, n2 := range n.children {
		n2.print(method)
	}
}

func (r *Route) printRoute() {
	for s, n := range r.trie {
		n.print(s)
	}
}

func (n *node) search(key string) string {
	
}

func (r *Route) getRoute(ctx *Context) {
	for s, n := range r.trie {

	}
}
