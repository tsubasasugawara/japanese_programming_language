package ast

type Program struct {
	Nodes []*Node
}

func NewProgram() *Program {
	p := &Program{}
	return p
}
