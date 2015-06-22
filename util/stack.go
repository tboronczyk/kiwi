package util

type Stack []interface{}

func NewStack() Stack {
	return Stack{}
}

func (s *Stack) Push(e interface{}) {
	(*s) = append((*s), e)
}

func (s Stack) Peek() interface{} {
	if s.Size() == 0 {
		panic("Attempt to access an empty stack")
	}
	return s[s.Size()-1]
}

func (s *Stack) Pop() interface{} {
	defer func() {
		(*s) = (*s)[:s.Size()-1]
	}()
	return s.Peek()
}

func (s Stack) Size() int {
	return len(s)
}
