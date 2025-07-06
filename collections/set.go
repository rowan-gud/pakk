package collections

type Set[T comparable] struct {
	data map[T]bool
	len  int
}

func NewSet[T comparable](elems ...T) *Set[T] {
	l := len(elems)
	m := make(map[T]bool, l)

	for _, elem := range elems {
		m[elem] = true
	}

	return &Set[T]{
		data: m,
		len:  l,
	}
}

func (s *Set[T]) Add(elem T) {
	had := s.data[elem]
	s.data[elem] = true

	if !had {
		s.len += 1
	}
}

func (s *Set[T]) Clone() *Set[T] {
	m := make(map[T]bool, s.len)

	for k := range s.data {
		m[k] = true
	}

	return &Set[T]{
		data: m,
		len:  s.len,
	}
}

func (s *Set[T]) Delete(elem T) bool {
	had := s.data[elem]

	delete(s.data, elem)

	if had {
		s.len -= 1
	}

	return had
}

func (s *Set[T]) Empty() bool {
	return s.len == 0
}

func (s *Set[T]) Has(elem T) bool {
	return s.data[elem]
}

func (s *Set[T]) Len() int {
	return s.len
}
