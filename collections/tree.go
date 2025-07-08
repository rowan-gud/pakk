package collections

type Tree[T any] struct {
	data     T
	children []*Tree[T]
}

func NewTree[T any](data T, children ...*Tree[T]) *Tree[T] {
	return &Tree[T]{
		data:     data,
		children: children,
	}
}

func (n *Tree[T]) AddChild(child *Tree[T]) {
	n.children = append(n.children, child)
}

func (n *Tree[T]) Children() []*Tree[T] {
	return n.children
}

func (n *Tree[T]) Data() T {
	return n.data
}

func (n *Tree[T]) Height() int {
	h := 1

	for _, child := range n.children {
		h = max(h, child.Height()+1)
	}

	return h
}
