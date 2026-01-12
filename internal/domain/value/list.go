package value

type List struct {
	V []string
}

func (List) Type() Type {
	return TypeList
}
