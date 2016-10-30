package router

import (
	"fmt"
	"strings"
)

// node is vertex of trie
type node struct {
	edges map[byte]*node
	value HandlerFunc
	path  string
	param *node
}

const (
	capacity    uint8 = 30
	paramSymbol byte  = ':'
)

const (
	static uint8 = iota // default
	param
)

/*
newNode returns a new initialezed *node
*/
func newNode(path string, handler HandlerFunc) *node {
	return &node{
		path:  path,
		value: handler,
		edges: make(map[byte]*node, capacity),
		param: nil,
	}
}

func newRootNode() *node {
	return newNode("", nil)
}

func (n *node) reset(path string) {
	n.path = path
	n.value = nil
	n.edges = make(map[byte]*node, capacity)
	n.param = nil
}

// Insert new path to node, if it isn't exist
func (n *node) insert(path string, handler HandlerFunc) {
	if len(path) < 1 {
		panic("Empty path can't be inserted")
	}

	// Count longest common prefix
	lcp := lcp(n.path, path)

	// path equal n.path
	if lcp == len(n.path) && lcp == len(path) {
		if n.value != nil {
			panic("Path already exist in trie.")
		}

		n.value = handler

		return
	}

	// Split
	if lcp > 0 && lcp < len(n.path) {
		n.split(lcp)

		n.insert(path, handler)

		return
	}

	path = path[lcp:]

	firstChar := path[0]

	edge := n.edges[firstChar]

	// if prefix doesn't exist in trie create new node represent this prefix
	if edge == nil {
		// Check if path has params
		paramIndex, paramLength := findFirstParam(path, paramSymbol)

		if paramIndex == -1 {
			if n.param != nil {
				panic("It's not possible to add static route with wildcard!")
			}

			n.edges[firstChar] = newNode(path, handler)

			return
		}

		// If param name is empty
		if paramLength < 2 {
			panic(fmt.Sprintf("Empty param name in path \"%s\"", path))
		}

		paramNode := newNode(path[paramIndex+1:paramIndex+paramLength], handler)

		if paramIndex+paramLength < len(path) {
			paramNode.value = nil

			paramNode.insert(path[paramIndex+paramLength:], handler)
		}

		if paramIndex == 0 {
			if firstChar != paramSymbol {
				panic("It's not possible to add static route with wildcard!")
			}

			if n.param != nil {
				if len(path) <= paramIndex+paramLength {
					panic("Path already exist!")
				}

				n.param.insert(path[paramIndex+paramLength:], handler)

				return
			}

			n.param = paramNode

			return
		}

		newNode := newNode(path[:paramIndex], nil)
		newNode.param = paramNode

		n.edges[firstChar] = newNode

		return
	}

	edge.insert(path, handler)
}

func (n *node) split(length int) {
	splittedNode := &node{
		path:  n.path[length:],
		value: n.value,
		edges: n.edges,
		param: n.param,
	}

	n.reset(n.path[:length])

	n.edges[splittedNode.path[0]] = splittedNode
}

func (n *node) lookup(path string) (handler HandlerFunc, params Params) {
	if len(path) < 1 {
		panic("Empty path can't be found")
	}

	for n != nil {
		lcp := lcp(n.path, path)

		if lcp < len(n.path) {
			return nil, params
		}

		if lcp == len(path) {
			return n.value, params
		}

		path = path[lcp:]

		nextNode := n.edges[path[0]]

		if nextNode == nil && n.param != nil {
			separatorIndex := strings.IndexByte(path, '/')
			if separatorIndex == -1 {
				separatorIndex = len(path)
			}

			paramValue := path[:separatorIndex]
			paramName := n.param.path

			if params == nil {
				params = make(Params)
			}

			params[paramName] = paramValue

			if separatorIndex == len(path) {
				return n.param.value, params
			}

			path = path[separatorIndex:]

			nextNode = n.param.edges[path[0]]
		}

		n = nextNode
	}

	return nil, params
}
