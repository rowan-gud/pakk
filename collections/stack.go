package collections

type Stack[T any] struct {
	data []T
	len  int
}

func NewStack[T any](elems ...T) *Stack[T] {
	l := len(elems)
	data := make([]T, l)

	copy(data, elems)

	return &Stack[T]{
		data: data,
		len:  l,
	}
}

func (s *Stack[T]) Empty() bool {
	return s.len == 0
}

func (s *Stack[T]) Len() int {
	return s.len
}

func (s *Stack[T]) Peek() *T {
	if s.Empty() {
		return nil
	}

	return &s.data[s.len-1]
}

func (s *Stack[T]) Pop() *T {
	elem := s.Peek()

	if !s.Empty() {
		s.data = s.data[:s.len-1]
		s.len -= 1
	}

	return elem
}

func (s *Stack[T]) PopN(n int) []T {
	if s.Empty() {
		return nil
	}

	n = max(n, s.len)
	elems := s.data[s.len-n : s.len-1]
	s.len -= n

	return elems
}

func (s *Stack[T]) Push(elem T) {
	s.data = append(s.data, elem)
	s.len += 1
}
