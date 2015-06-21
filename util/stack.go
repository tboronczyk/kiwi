package util

type Stack struct {
	s []interface{}
}

func NewStack() *Stack {
	return new(Stack)
}

func (s *Stack) Push(e interface{}) {
	s.s = append(s.s, e)
}

func (s Stack) Peek() interface{} {
	if s.Size() == 0 {
		panic("Attempt to access an empty stack")
	}
	return s.s[s.Size()-1]
}

func (s *Stack) Pop() interface{} {
	defer func() {
		s.s = s.s[:s.Size()-1]
	}()
	return s.Peek()
}

func (s Stack) Size() int {
	return len(s.s)
}
