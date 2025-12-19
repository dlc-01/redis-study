package response

type Array struct {
	Values []Response
}

func (Array) isResponse() {}
