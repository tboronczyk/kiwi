package util

// Stack is a standard stack data structure.
type Stack []interface{}

// NewStack returns an implementation of a stack.
func NewStack() Stack {
	return make(Stack, 0)
}

// Push places entry e on top of the stack.
func (s *Stack) Push(e interface{}) {
	(*s) = append((*s), e)
}

// Peek returns the entry on top of the stack.
func (s Stack) Peek() interface{} {
	if s.Size() == 0 {
		panic("Attempt to access an empty stack")
	}
	return s[s.Size()-1]
}

// Push removes and returns the entry from the top of the stack.
func (s *Stack) Pop() interface{} {
	defer func() {
		(*s) = (*s)[:s.Size()-1]
	}()
	return s.Peek()
}

// Size returns the number of entries on the stack.
func (s Stack) Size() int {
	return len(s)
}
