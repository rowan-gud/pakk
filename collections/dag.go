package collections

type dagNode[K comparable, T any] struct {
	key    K
	data   T
	edges  map[K]*dagNode[K, T]
	isRoot bool
	isLeaf bool
}

type DAG[K comparable, T any] struct {
	nodes map[K]*dagNode[K, T]
}

func NewDAG[K comparable, T any]() *DAG[K, T] {
	return &DAG[K, T]{}
}

func (g *DAG[K, T]) AddNode(key K, data T) int {
	g.nodes[key] = &dagNode[K, T]{
		key:    key,
		data:   data,
		isLeaf: true,
		isRoot: true,
	}

	return len(g.nodes) - 1
}

func (g *DAG[K, T]) AddEdge(from, to K) {
	fromNode, fromExists := g.nodes[from]
	toNode, toExists := g.nodes[to]

	if !fromExists || !toExists {
		return
	}

	fromNode.isLeaf = false
	toNode.isRoot = false

	fromNode.edges[to] = toNode
}
