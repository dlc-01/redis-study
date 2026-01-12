package value

type String struct {
	V string
}

func (String) Type() Type {
	return TypeString
}
