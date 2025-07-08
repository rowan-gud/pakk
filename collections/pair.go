package collections

type Pair[T1, T2 any] struct {
	A T1
	B T2
}

func NewPair[T1, T2 any](a T1, b T2) Pair[T1, T2] {
	return Pair[T1, T2]{
		A: a,
		B: b,
	}
}
