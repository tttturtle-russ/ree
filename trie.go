package ree

import (
	rlog "Ree/ree/logs"
	"strings"
)

type node struct {
	isWild   bool             // 是否是*
	children map[string]*node // 子节点
	part     string           // /分割的路径
	path     string           // 完整的路径
}

func (n *node) insert(path string, parts []string, height int) {
	if len(parts) == height {
		n.path = path
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{
			isWild:   part[0] == ':' || part[0] == '*',
			children: make(map[string]*node),
			part:     part,
		}
		n.children[part] = child
	}
	child.insert(path, parts, height+1)
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.path == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	children := n.matchChildren(part)
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}

func (r *Route) getRoute(method, path string) (*node, map[string]string) {
	parts := parsePath(path)
	params := make(map[string]string)
	if _, ok := r.trie[method]; !ok {
		return nil, nil
	}
	n := r.trie[method].search(parts, 0)
	if n != nil {
		nParts := parsePath(n.path)
		for i, part := range nParts {
			if part[0] == ':' {
				params[part[1:]] = parts[i]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(nParts[i:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

func parsePath(path string) []string {
	tmp := strings.Split(path, "/")
	parts := make([]string, 0)
	for _, part := range tmp {
		if part != "" {
			parts = append(parts, part)
			if part[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *Route) addRoute(method string, path string, handler HandlerFunc) {
	parts := parsePath(path)
	if _, ok := r.trie[method]; !ok {
		r.trie[method] = &node{children: make(map[string]*node)}
	}
	n := r.trie[method]
	n.insert(path, parts, 0)
	r.handlers[method+"-"+path] = handler
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
