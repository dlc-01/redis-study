package response

type SimpleString struct {
	Value string
}

func (SimpleString) isResponse() {}
