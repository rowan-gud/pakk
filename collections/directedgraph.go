package collections

import "github.com/BurntSushi/toml"

type DirectedGraphNode[K comparable, T any] struct {
	Key    K
	data   T
	edges  map[K]*DirectedGraphNode[K, T]
	IsLeaf bool
	IsRoot bool
}

type DirectedGraph[K comparable, T any] struct {
	nodes map[K]*DirectedGraphNode[K, T]
}

func NewDirectedGraph[K comparable, T any]() *DirectedGraph[K, T] {
	return &DirectedGraph[K, T]{
		nodes: make(map[K]*DirectedGraphNode[K, T]),
	}
}

func (g *DirectedGraph[K, T]) AddNode(key K, data T) {
	g.nodes[key] = &DirectedGraphNode[K, T]{
		Key:    key,
		data:   data,
		edges:  make(map[K]*DirectedGraphNode[K, T]),
		IsLeaf: true,
		IsRoot: true,
	}
}

func (g *DirectedGraph[K, T]) AddEdge(from K, to K) {
	a, exists := g.nodes[from]
	if !exists {
		return
	}

	a.IsLeaf = false

	b, exists := g.nodes[to]
	if !exists {
		return
	}

	b.IsRoot = false

	a.edges[to] = b
}

func (g *DirectedGraph[K, T]) Node(key K) (*DirectedGraphNode[K, T], bool) {
	n, exists := g.nodes[key]

	return n, exists
}

func (g *DirectedGraph[K, T]) Nodes() map[K]*DirectedGraphNode[K, T] {
	return g.nodes
}

func (g *DirectedGraph[K, T]) MarshalTOML() ([]byte, error) {
	nodes := []map[string]any{}
	edges := []map[string]K{}

	for _, n := range g.nodes {
		nodes = append(nodes, map[string]any{
			"is_leaf": n.IsLeaf,
			"is_root": n.IsRoot,
			"data":    n.data,
		})

		for to := range n.edges {
			edges = append(edges, map[string]K{
				"from": n.Key,
				"to":   to,
			})
		}
	}

	return toml.Marshal(map[string]any{
		"nodes": nodes,
		"edges": edges,
	})
}

func (n *DirectedGraphNode[K, T]) Data() T {
	return n.data
}

func (n *DirectedGraphNode[K, T]) Edge(key K) (*DirectedGraphNode[K, T], bool) {
	edge, exists := n.edges[key]

	return edge, exists
}

func (n *DirectedGraphNode[K, T]) Edges() []*DirectedGraphNode[K, T] {
	res := make([]*DirectedGraphNode[K, T], len(n.edges))

	idx := 0

	for _, edge := range n.edges {
		res[idx] = edge

		idx += 1
	}

	return res
}
