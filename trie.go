// 前缀树路由实现
package gee

import (
	"fmt"
	"strings"
)

type (
	Trie struct {
		root *node
	}
	node struct {
		// 匹配的url(注册时提供的url，例如：/p/:lang/doc)
		pattern string
		// url中的一段，例如：p  :lang  doc
		part string
		// 子节点
		children []*node
		// 任意节点都能匹配，例如当前节点的part是:lang或者*filepath，则IsWild为true
		isWild bool
	}
)

func NewTrie() *Trie {
	root := &node{
		children: make([]*node, 0),
	}
	return &Trie{
		root: root,
	}
}

func (t Trie) parsePattern(pattern string) []string {
	parts := strings.Split(pattern, "/")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if part == "" || part == "/" {
			continue
		}
		result = append(result, part)
	}
	return result
}

func (t *Trie) matchChild(n *node, part string) *node {
	for _, child := range n.children {
		// 找到匹配的part
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

func (t *Trie) matchChildren(n *node, part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		// 找到匹配的part
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (t *Trie) Insert(pattern string) {
	parts := t.parsePattern(pattern)
	found, _ := t.Search(pattern)
	if found != nil {
		panic(fmt.Sprintf("Repeated pattern between %s and %s ", pattern, found.pattern))
	}
	t.insert(t.root, pattern, parts, 0)
}

func (t *Trie) insert(n *node, pattern string, parts []string, index int) {
	// 最后一part，打上pattern
	if len(parts) == index {
		n.pattern = pattern
		return
	}

	part := parts[index]
	findNode := t.matchChild(n, part)
	// 没找到，创建一个新的节点
	if findNode == nil {
		findNode = new(node)
		findNode.part = part
		// *filepath/:param，参数可以匹配任意值
		findNode.isWild = part[0] == '*' || part[0] == ':'
		n.children = append(n.children, findNode)
	}

	t.insert(findNode, pattern, parts, index+1)
}

func (t Trie) Search(pattern string) (*node, map[string]string) {
	parts := t.parsePattern(pattern)
	n := t.search(t.root, pattern, parts, 0)
	var params map[string]string
	if n != nil {
		params = t.getParams(n.pattern, pattern)
	}
	return n, params
}

func (t Trie) search(n *node, pattern string, parts []string, index int) *node {
	// 当前节点是最后一个
	if len(parts) == index {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[index]
	// 当前节点不是最后一个，继续匹配子节点
	children := t.matchChildren(n, part)
	for _, child := range children {
		// *filepath
		if strings.HasPrefix(child.part, "*") {
			return child
		}
		found := t.search(child, pattern, parts, index+1)
		if found != nil {
			return found
		}
	}
	// 子节点中没找到匹配的节点
	return nil
}

// 从2个pattern中匹配参数
// pattern 已经注册的节点的pattern
// cgi 用户访问的url
func (t Trie) getParams(pattern, cgi string) map[string]string {
	params := make(map[string]string, 0)
	parts := t.parsePattern(pattern)
	cgis := t.parsePattern(cgi)

	for i, part := range parts {
		if part[0] == ':' {
			params[part[1:]] = cgis[i]
			continue
		}
		if part[0] == '*' {
			params[part[1:]] = strings.Join(cgis[i:], "/")
			break
		}
	}
	return params
}
