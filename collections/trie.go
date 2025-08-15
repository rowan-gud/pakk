package collections

import (
	"fmt"
	"strings"
)

type Trie[K comparable, T any] struct {
	root *trieNode[K, T]
}

func NewTrie[K comparable, T any]() *Trie[K, T] {
	return &Trie[K, T]{
		root: newTrieNode[K, T](),
	}
}

func (t *Trie[K, T]) Add(path []K, data *T) {
	cur := t.root

	for _, key := range path {
		cur = cur.child(key)
	}

	cur.data = data
}

func (t *Trie[K, T]) Get(path []K) (*T, error) {
	var err error
	cur := t.root

	for _, key := range path {
		cur, err = cur.requireChild(key)
		if err != nil {
			return nil, err
		}
	}

	return cur.data, nil
}

func (t *Trie[K, T]) Print() {
	fmt.Println("root:")
	t.root.print(0)
}

type trieNode[K comparable, T any] struct {
	data     *T
	children map[K]*trieNode[K, T]
}

func newTrieNode[K comparable, T any]() *trieNode[K, T] {
	return &trieNode[K, T]{
		children: make(map[K]*trieNode[K, T]),
	}
}

func (n *trieNode[K, T]) requireChild(key K) (*trieNode[K, T], error) {
	if n.children[key] == nil {
		return nil, fmt.Errorf("child %v not found", key)
	}

	return n.children[key], nil
}

func (n *trieNode[K, T]) child(key K) *trieNode[K, T] {
	if n.children[key] == nil {
		n.children[key] = newTrieNode[K, T]()
	}

	return n.children[key]
}

func (n *trieNode[K, T]) print(indent int) {
	for key, child := range n.children {
		fmt.Println(strings.Repeat(" ", indent), key, child.data)
		child.print(indent + 2)
	}
}
