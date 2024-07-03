package trie

import (
	"strings"
)

type Trie struct {
	root       *Node
	singleWild string
	multiWild  string
	separator  string
	isSWild    bool
	isMWild    bool
}

type Node struct {
	children map[string]*Node
	value    interface{}
}

type Option func(*Trie)

type KeyValue struct {
	Key   string
	Value interface{}
}

type stackNode struct {
	node    *Node
	i       int
	keypart string
}

type trieError struct {
	message string
}

func (e trieError) Error() string {
	return e.message
}

func WithSingleWild(wild string) Option {
	return func(t *Trie) {
		t.singleWild = wild
		t.isSWild = true
	}
}

func WithMultiWild(wild string) Option {
	return func(t *Trie) {
		t.multiWild = wild
		t.isMWild = true
	}
}

func WithSeparator(sep string) Option {
	return func(t *Trie) {
		t.separator = sep
	}
}

func New(options ...Option) *Trie {
	t := &Trie{
		root:      &Node{children: make(map[string]*Node)},
		separator: "", // Default separator is an empty string
	}
	for _, opt := range options {
		opt(t)
	}
	return t
}

func (t *Trie) Insert(key string, value interface{}) (bool, error) {
	node := t.root
	parts := strings.Split(key, t.separator)
	for i, part := range parts {
		if t.isMWild && part == t.multiWild && i != len(parts)-1 {
			return false, trieError{message: "error: multi-level wildcard '" + t.multiWild + "' permitted only at the end of the insert key."}
		}
		if _, exists := node.children[part]; !exists {
			node.children[part] = &Node{children: make(map[string]*Node)}
		}
		node = node.children[part]
	}
	node.value = value
	return true, nil
}

func (t *Trie) Retrieve(key string) (interface{}, error) {
	node := t.root
	parts := strings.Split(key, t.separator)
	for i, part := range parts {
		if t.isMWild && part == t.multiWild && i != len(parts)-1 {
			return nil, trieError{message: "error: multi-level wildcard '" + t.multiWild + "' permitted only at the end of the retrieve key."}
		}
		nextNode, exists := node.children[part]
		if !exists {
			return nil, nil
		}
		node = nextNode
	}
	return node.value, nil
}

func collectAll(startNode *Node, startKeypart string, matches *[]KeyValue, separator string) {
	stack := []stackNode{{node: startNode, keypart: startKeypart}}
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if current.node.value != nil {
			*matches = append(*matches, KeyValue{Key: current.keypart, Value: current.node.value})
		}
		for k, v := range current.node.children {
			stack = append(stack, stackNode{node: v, keypart: current.keypart + k + separator})
		}
	}
}

func (t *Trie) Match(key string) []KeyValue {
	var matches []KeyValue

	parts := strings.Split(key, t.separator)
	stack := []stackNode{{node: t.root, i: 0, keypart: ""}}

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		node, i, keypart := current.node, current.i, current.keypart

		if t.isMWild && parts[i] == t.multiWild {
			collectAll(node, keypart, &matches, t.separator)
		} else if t.isSWild && parts[i] == t.singleWild {
			for k, childNode := range node.children {
				if i == len(parts)-1 && childNode.value != nil {
					matches = append(matches, KeyValue{Key: keypart + k, Value: childNode.value})
				} else if i < len(parts)-1 {
					stack = append(stack, stackNode{node: childNode, i: i + 1, keypart: keypart + k + t.separator})
				}
			}
		} else {
			keyNode, exists := node.children[parts[i]]
			if exists {
				if i == len(parts)-1 && keyNode.value != nil {
					matches = append(matches, KeyValue{Key: keypart + parts[i], Value: keyNode.value})
				} else {
					stack = append(stack, stackNode{node: keyNode, i: i + 1, keypart: keypart + parts[i] + t.separator})
				}
			}

			if t.isSWild {
				singleWildNode := node.children[t.singleWild]
				if singleWildNode != nil {
					if i == len(parts)-1 && singleWildNode.value != nil {
						matches = append(matches, KeyValue{Key: keypart + t.singleWild, Value: singleWildNode.value})
					} else {
						stack = append(stack, stackNode{node: singleWildNode, i: i + 1, keypart: keypart + t.singleWild + t.separator})
					}
				}
			}

			if t.isMWild {
				multiWildNode := node.children[t.multiWild]
				if multiWildNode != nil {
					matches = append(matches, KeyValue{Key: keypart + t.multiWild, Value: multiWildNode.value})
				}
			}

		}
	}
	return matches
}

func (t *Trie) Delete(key string) bool {
	parentStack := []*Node{t.root}
	node := t.root
	parts := strings.Split(key, t.separator)
	for _, part := range parts {
		child, exists := node.children[part]
		if !exists {
			return false
		}
		parentStack = append(parentStack, child)
		node = child
	}
	if node.value == nil {
		return false
	}
	node.value = nil
	for i := len(parts) - 1; i >= 0; i-- {
		if node.value != nil || len(node.children) > 0 {
			break
		}
		parentStack = parentStack[:len(parentStack)-1]
		parent := parentStack[len(parentStack)-1]
		delete(parent.children, parts[i])
		node = parent
	}
	return true
}
