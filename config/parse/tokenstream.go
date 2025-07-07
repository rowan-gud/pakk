package parse

import "strings"

type tokenStream struct {
	chars []string
	curr  int
	len   int
}

func newTokenStream(str string) *tokenStream {
	chars := strings.Split(str, "")

	return &tokenStream{
		chars: chars,
		curr:  0,
		len:   len(chars),
	}
}

func (s *tokenStream) Next() bool {
	return s.curr < s.len
}

func (s *tokenStream) Peek() string {
	if s.Next() {
		return s.chars[s.curr]
	}

	return ""
}

func (s *tokenStream) Reset() {
	s.Seek(0)
}

func (s *tokenStream) Seek(pos int) {
	s.curr = pos
}

func (s *tokenStream) Take() string {
	tok := s.chars[s.curr]
	s.curr += 1

	return tok
}

func (s *tokenStream) TakeWhile(fn func(token string) bool) string {
	res := ""

	for s.Next() {
		tok := s.Peek()

		if fn(tok) {
			res += s.Take()
		} else {
			break
		}
	}

	return res
}
